package tasks

import (
	"context"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rcommon "github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
	"github.com/jackc/pgx/v4"
)

type QueryerWrap common.QueryerWrap

func (w QueryerWrap) GetAndBlockIdsByColumn(columnId rcommon.Id) ([]rcommon.Id, error) {
	var ids = []rcommon.Id{}
	const q = `SELECT id FROM tasks WHERE column_id = $1 ORDER BY rank ASC FOR UPDATE`
	rows, err := w.Q.Query(context.Background(), q, columnId)
	if err != nil {
		return ids, err
	}
	defer rows.Close()
	var id rcommon.Id
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (w QueryerWrap) GetAndBlockMaxRankByColumn(columnId rcommon.Id) (rank rcommon.Rank, err error) {
	const q = `
		SELECT rank FROM tasks
		WHERE column_id = $1
		ORDER BY rank DESC
		LIMIT 1
	`
	err = w.Q.QueryRow(context.Background(), q, columnId).Scan(&rank)
	return rank, err
}

func (w QueryerWrap) UpdatePosition(taskId, columnId rcommon.Id, rank rcommon.Rank) error {
	const q = "UPDATE tasks SET column_id = $2, rank = $3 WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, taskId, columnId, rank))
}

func (w QueryerWrap) Get(taskId rcommon.Id) (rcommon.Task, error) {
	t := rcommon.Task{Id: taskId}
	const q = `SELECT project_id, column_id, name, description FROM tasks WHERE id = $1`
	err := w.Q.QueryRow(context.Background(), q, taskId).Scan(&t.ProjectId, &t.ColumnId, &t.Name, &t.Description)
	return t, err
}

func (w QueryerWrap) GetExpanded(taskId rcommon.Id) (rcommon.TaskExpanded, error) {
	const q = `
		SELECT t.id, t.project_id, t.column_id, t.name, t.description, 
			   COALESCE(c.id, 0), COALESCE(c.text, '')
		FROM tasks t
		LEFT JOIN comments c ON c.task_id = t.id
		WHERE t.id = $1
		ORDER BY c.create_dt ASC
	`
	rows, err := w.Q.Query(context.Background(), q, taskId)
	if err != nil {
		return rcommon.TaskExpanded{}, err
	}
	defer rows.Close()
	return buildExpanded(rows)
}

func buildExpanded(rows pgx.Rows) (rcommon.TaskExpanded, error) {
	t := rcommon.Task{}
	c := rcommon.Comment{}
	comments := []rcommon.Comment{}
	for rows.Next() {
		err := rows.Scan(&t.Id, &t.ProjectId, &t.ColumnId, &t.Name, &t.Description, &c.Id, &c.Text)
		if err != nil {
			return rcommon.TaskExpanded{}, err
		}
		if c.Id != 0 {
			comments = append(comments, c)
		}
	}
	if t.Id != 0 {
		return rcommon.TaskExpanded{Task: t, Comments: comments}, nil
	} else {
		return rcommon.TaskExpanded{}, common.ErrNoRows
	}
}

func (w QueryerWrap) Create(projectId, columnId rcommon.Id, name string,
	description string, rank rcommon.Rank) (rcommon.Task, error) {
	t := rcommon.Task{
		ProjectId: projectId, ColumnId: columnId,
		TaskSettableFields: rcommon.TaskSettableFields{Name: name, Description: description},
	}
	const q = `
		INSERT INTO tasks (project_id, column_id, name, description, rank) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := w.Q.QueryRow(context.Background(), q, projectId, columnId, name, description, rank).Scan(&t.Id)
	return t, err
}

func (w QueryerWrap) GetAndBlockRank(columnId, taskId rcommon.Id) (rank rcommon.Rank, err error) {
	const q = `SELECT rank FROM tasks WHERE column_id = $1 AND id = $2 FOR UPDATE`
	err = w.Q.QueryRow(context.Background(), q, columnId, taskId).Scan(&rank)
	return rank, err
}

func (w QueryerWrap) GetNextRank(columnId rcommon.Id, rank rcommon.Rank) (nextRank rcommon.Rank, err error) {
	const q = `
		SELECT rank FROM tasks
		WHERE column_id = $1 AND rank > $2
		ORDER BY rank
		LIMIT 1
	`
	err = w.Q.QueryRow(context.Background(), q, columnId, rank).Scan(&nextRank)
	return nextRank, err
}

func (w QueryerWrap) Update(taskId rcommon.Id, name string, description string) error {
	const q = "UPDATE tasks SET name = $2, description = $3 WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, taskId, name, description))
}

func (w QueryerWrap) Delete(taskId rcommon.Id) error {
	const q = "DELETE FROM tasks WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, taskId))
}
