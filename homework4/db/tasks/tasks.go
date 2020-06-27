package tasks

import (
	"context"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
	"github.com/jackc/pgx/v4"
)

type QueryerWrap common.QueryerWrap

func (w QueryerWrap) GetAndBlockIdsByColumn(columnId rsrc.Id) (ids []rsrc.Id, err error) {
	const q = `SELECT id FROM tasks WHERE column_id = $1 ORDER BY rank ASC`
	rows, err := w.Q.Query(context.Background(), q, columnId)
	if err != nil {
		return ids, err
	}
	defer rows.Close()
	var id rsrc.Id
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (w QueryerWrap) GetAndBlockMaxRankByColumn(columnId rsrc.Id) (rank rsrc.Rank, err error) {
	const q = `
		SELECT rank FROM tasks
		WHERE column_id = $1
		ORDER BY rank DESC
		LIMIT 1
	`
	err = w.Q.QueryRow(context.Background(), q, columnId).Scan(&rank)
	return rank, err
}

func (w QueryerWrap) UpdatePosition(taskId, columnId rsrc.Id, rank rsrc.Rank) error {
	const q = "UPDATE tasks SET column_id = $2, rank = $3 WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, taskId, columnId, rank))
}

func (w QueryerWrap) Get(taskId rsrc.Id) (rsrc.Task, error) {
	t := rsrc.Task{Id: taskId}
	const q = `SELECT project_id, column_id, name, description FROM tasks WHERE id = $1`
	err := w.Q.QueryRow(context.Background(), q, taskId).Scan(&t.ProjectId, &t.ColumnId, &t.Name, &t.Description)
	return t, err
}

func (w QueryerWrap) GetExpanded(taskId rsrc.Id) (rsrc.TaskExpanded, error) {
	const q = `
		SELECT t.project_id, t.column_id, t.name, t.description, c.id, c.text
		FROM tasks t
		JOIN comments c
		WHERE id = $1
		ORDER BY c.create_dt ASC
	`
	rows, err := w.Q.Query(context.Background(), q, taskId)
	if err != nil {
		return rsrc.TaskExpanded{}, err
	}
	defer rows.Close()
	return buildExpanded(taskId, rows)
}

func buildExpanded(taskId rsrc.Id, rows pgx.Rows) (rsrc.TaskExpanded, error) {
	t := rsrc.Task{Id: taskId}
	c := rsrc.Comment{}
	comments := []rsrc.Comment{}
	for rows.Next() {
		err := rows.Scan(&t.ProjectId, t.ColumnId, t.Name, t.Description, c.Id, c.Text)
		if err != nil {
			return rsrc.TaskExpanded{}, err
		}
		comments = append(comments, c)
	}
	if len(comments) > 0 {
		return rsrc.TaskExpanded{Task: t, Comments: comments}, nil
	} else {
		return rsrc.TaskExpanded{}, common.ErrNoRows
	}
}

func (w QueryerWrap) Create(projectId, columnId rsrc.Id, name string,
	description string, rank rsrc.Rank) (rsrc.Task, error) {
	t := rsrc.Task{
		ProjectId: projectId, ColumnId: columnId,
		TaskSettableFields: rsrc.TaskSettableFields{Name: name, Description: description},
	}
	const q = `
		INSERT INTO tasks (project_id, column_id, name, description, rank) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := w.Q.QueryRow(context.Background(), q, projectId, columnId, name, description, rank).Scan(&t.Id)
	return t, err
}

func (w QueryerWrap) GetAndBlockRank(columnId, taskId rsrc.Id) (rank rsrc.Rank, err error) {
	const q = `SELECT rank FROM tasks WHERE columnId = $1 AND taskId = $2 FOR UPDATE`
	err = w.Q.QueryRow(context.Background(), q, columnId, taskId).Scan(&rank)
	return rank, err
}

func (w QueryerWrap) GetNextRank(columnId rsrc.Id, rank rsrc.Rank) (nextRank rsrc.Rank, err error) {
	const q = `
		SELECT rank FROM tasks
		WHERE columnId = $1 AND rank > $2
		ORDER BY rank
		LIMIT 1
	`
	err = w.Q.QueryRow(context.Background(), q, columnId, rank).Scan(&nextRank)
	return nextRank, err
}

func (w QueryerWrap) Update(taskId rsrc.Id, name string, description string) error {
	const q = "UPDATE tasks SET name = $2, description = $3 WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, taskId, name, description))
}

func (w QueryerWrap) Delete(taskId rsrc.Id) error {
	const q = "DELETE FROM tasks WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, taskId))
}
