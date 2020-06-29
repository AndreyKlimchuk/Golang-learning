package columns

import (
	"context"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rcommon "github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
)

type QueryerWrap common.QueryerWrap

func (w QueryerWrap) GetAndBlockMaxRank(projectId rcommon.Id) (rank rcommon.Rank, err error) {
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

func (w QueryerWrap) Create(projectId rcommon.Id, name string, rank rcommon.Rank) (rcommon.ColumnExpanded, error) {
	c := rcommon.ColumnExpanded{
		Column: rcommon.Column{ColumnSettableFields: rcommon.ColumnSettableFields{Name: name}},
		Tasks:  []rcommon.Task{},
	}
	const q = `INSERT INTO columns (project_id, name, rank) VALUES ($1, $2, $3) RETURNING id`
	err := w.Q.QueryRow(context.Background(), q, projectId, name, rank).Scan(&c.Id)
	return c, err
}

func (w QueryerWrap) Get(projectId, columnId rcommon.Id) (rcommon.Column, error) {
	c := rcommon.Column{Id: columnId}
	const q = `SELECT name FROM columns WHERE project_id = $1 AND id = $2`
	err := w.Q.QueryRow(context.Background(), q, projectId, columnId).Scan(&c.Name)
	return c, err
}

func (w QueryerWrap) GetMultiple(projectId rcommon.Id) ([]rcommon.Column, error) {
	columns := []rcommon.Column{}
	const q = "SELECT id, name FROM columns WHERE project_id = $1 ORDER BY rank ASC"
	rows, err := w.Q.Query(context.Background(), q, projectId)
	if err != nil {
		return columns, err
	}
	defer rows.Close()
	c := rcommon.Column{}
	for rows.Next() {
		err := rows.Scan(&c.Id, &c.Name)
		if err != nil {
			return columns, err
		}
		columns = append(columns, c)
	}
	return columns, nil
}

func (w QueryerWrap) Update(projectId, columnId rcommon.Id, name string) error {
	const q = `UPDATE columns SET name = $3 WHERE project_id = $1 AND id = $2`
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, projectId, columnId, name))
}

func (w QueryerWrap) GetAndBlockRank(projectId, columnId rcommon.Id) (rank rcommon.Rank, err error) {
	const q = `SELECT rank FROM columns WHERE project_id = $1 AND id = $2 FOR UPDATE`
	err = w.Q.QueryRow(context.Background(), q, projectId, columnId).Scan(&rank)
	return rank, err
}

func (w QueryerWrap) Delete(columnId rcommon.Id) error {
	const q = `DELETE FROM columns WHERE id = $1`
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, columnId))
}

func (w QueryerWrap) GetAndBlockSuccessorColumnId(projectId rcommon.Id, rank rcommon.Rank) (id rcommon.Id, err error) {
	const q = `
		WITH before AS (
			SELECT id FROM columns
			WHERE project_id = $1 AND rank < $2
			ORDER BY rank
			LIMIT 1
			FOR UPDATE
		), after AS (
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

func (w QueryerWrap) UpdateRank(projectId, columnId rcommon.Id, rank rcommon.Rank) error {
	const q = `UPDATE columns SET rank = $3 WHERE project_id = $1 AND id = $2`
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, projectId, columnId, rank))
}

func (w QueryerWrap) GetNextRank(projectId rcommon.Id, rank rcommon.Rank) (nextRank rcommon.Rank, err error) {
	const q = `
		SELECT rank FROM columns
		WHERE project_id = $1 AND rank > $2
		ORDER BY rank
		LIMIT 1
	`
	err = w.Q.QueryRow(context.Background(), q, projectId, rank).Scan(&nextRank)
	return nextRank, err
}

func (w QueryerWrap) GetByName(projectId rcommon.Id, name string) (rcommon.Column, error) {
	c := rcommon.Column{ColumnSettableFields: rcommon.ColumnSettableFields{Name: name}}
	const q = `SELECT id FROM columns WHERE project_id = $1 AND name = $2`
	err := w.Q.QueryRow(context.Background(), q, projectId, name).Scan(&c.Id)
	return c, err
}
