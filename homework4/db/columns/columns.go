package columns

import (
	"context"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type QueryerWrap common.QueryerWrap

func (w QueryerWrap) GetAndBlockMaxRank(projectId rsrc.Id) (rank rsrc.Rank, err error) {
	const q = `
		SELECT rank FROM columns
		WHERE project_id = $1
		ORDER BY rank DESC
		LIMIT 1
		FOR UPDATE
	`
	err = w.Q.QueryRow(context.Background(), q, projectId).Scan(&rank)
	return rank, err
}

func (w QueryerWrap) Create(projectId rsrc.Id, name string, rank rsrc.Rank) (rsrc.ColumnExpanded, error) {
	c := rsrc.ColumnExpanded{
		Column: rsrc.Column{ColumnSettableFields: rsrc.ColumnSettableFields{Name: name}},
		Tasks:  []rsrc.Task{},
	}
	const q = `INSERT INTO columns (project_id, name, rank) VALUES ($1, $2, $3)`
	err := w.Q.QueryRow(context.Background(), q, projectId, name, rank).Scan(&c.Id)
	return c, err
}

func (w QueryerWrap) Get(projectId, columnId rsrc.Id) (rsrc.Column, error) {
	c := rsrc.Column{Id: columnId}
	const q = `SELECT name FROM columns WHERE project_id = $1 AND id = $2`
	err := w.Q.QueryRow(context.Background(), q, projectId, columnId).Scan(&c.Name)
	return c, err
}

func (w QueryerWrap) GetMultiple(projectId rsrc.Id) (columns []rsrc.Column, err error) {
	const q = "SELECT id, name FROM columns WHERE project_id = $1 ORDER BY rank DESC"
	rows, err := w.Q.Query(context.Background(), q, projectId)
	if err != nil {
		return columns, err
	}
	defer rows.Close()
	c := rsrc.Column{}
	for rows.Next() {
		err := rows.Scan(c.Id, c.Name)
		if err != nil {
			return columns, err
		}
		columns = append(columns, c)
	}
	return columns, nil
}

func (w QueryerWrap) Update(projectId, columnId rsrc.Id, name string) error {
	const q = `UPDATE columns SET name = $3 WHERE project_id = $1 AND id = $2`
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, projectId, columnId, name))
}

func (w QueryerWrap) GetAndBlockRank(projectId, columnId rsrc.Id) (rank rsrc.Rank, err error) {
	const q = `SELECT rank FROM column WHERE project_id = $1 AND id = $2 FOR UPDATE`
	err = w.Q.QueryRow(context.Background(), q, projectId, columnId).Scan(&rank)
	return rank, err
}

func (w QueryerWrap) Delete(columnId rsrc.Id) error {
	const q = `DELETE FROM columns WHERE id = $1`
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, columnId))
}

func (w QueryerWrap) GetAndBlockSuccessorColumnId(projectId rsrc.Id, rank rsrc.Rank) (id rsrc.Id, err error) {
	const q = `
		WITH before AS (
			SELECT id FROM columns
			WHERE project_id = $1 AND rank < $2
			ORDER BY rank
			LIMIT 1
			FOR UPDATE
		WITH after AS (
			SELECT id FROM columns
			WHERE NOT EXISTS(SELECT * FROM before) AND project_id = $1 AND rank > $2
			ORDER BY rank
			LIMIT 1
			FOR UPDATE
		)
		SELECT * FROM before
		UNION
		SELECT * FROM after
	`
	err = w.Q.QueryRow(context.Background(), q, projectId, rank).Scan(&id)
	return id, err
}

func (w QueryerWrap) UpdateRank(projectId, columnId rsrc.Id, rank rsrc.Rank) error {
	const q = `UPDATE columns SET rank = $3 WHERE project_id = $1 AND id = $2`
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, projectId, columnId, rank))
}

func (w QueryerWrap) GetNextRank(projectId rsrc.Id, rank rsrc.Rank) (nextRank rsrc.Rank, err error) {
	const q = `
		SELECT rank FROM columns
		WHERE project_id = $1 AND rank > $2
		ORDER BY rank
		LIMIT 1
	`
	err = w.Q.QueryRow(context.Background(), q, projectId, rank).Scan(&nextRank)
	return nextRank, err
}

func (w QueryerWrap) GetByName(projectId rsrc.Id, name string) (rsrc.Column, error) {
	c := rsrc.Column{ColumnSettableFields: rsrc.ColumnSettableFields{Name: name}}
	const q = `SELECT id FROM columns WHERE project_id = $1 AND name = $2`
	err := w.Q.QueryRow(context.Background(), q, projectId, name).Scan(&c.Id)
	return c, err
}
