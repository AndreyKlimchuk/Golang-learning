package columns

import (
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type Query struct {
}

func (q Query) GetAndBlockMaxRank(projectId rsrc.Id) (rsrc.Rank, error) {
	return "", nil
}

func (q Query) Create(projectId rsrc.Id, name string, rank rsrc.Rank) (rsrc.ColumnExpanded, error) {
	return rsrc.ColumnExpanded{
		Column: rsrc.Column{
			ColumnSettableFields: rsrc.ColumnSettableFields{
				Name: name,
			},
		},
	}, nil
}

func (q Query) Get(projectId, columnId rsrc.Id) (rsrc.Column, error) {
	return rsrc.Column{Id: columnId}, nil
}

func (q Query) GetMultiple(projectId rsrc.Id) ([]rsrc.Column, error) {
	return []rsrc.Column{}, nil
}

func (q Query) Update(projectId, columnId rsrc.Id, name string) error {
	return nil
}

func (q Query) GetAndBlockRank(projectId, columnId rsrc.Id) (rsrc.Rank, error) {
	// sql := "DELETE FROM columns WHERE project_id = ? AND id = ? RETURNING rank"
	return "", nil
}

func (q Query) Delete(columnId rsrc.Id) error {
	// sql := "DELETE FROM columns WHERE project_id = ? AND id = ? RETURNING rank"
	return nil
}

func (q Query) GetAndBlockSuccessorColumnId(projectId rsrc.Id, r rsrc.Rank) (rsrc.Id, error) {
	// sql := "
	// 	WITH before AS (
	// 		SELECT id FROM columns
	// 		WHERE project_id = ? AND rank < ?
	// 		ORDER BY rank
	// 		LIMIT 1
	// 		FOR UPDATE
	// 	WITH after AS (
	// 		SELECT id FROM columns
	// 		WHERE NOT EXIST(SELECT * FROM before) AND project_id = ? AND rank > ?
	// 		ORDER BY rank
	// 		LIMIT 1
	// 		FOR UPDATE
	// 	)
	// 	SELECT * FROM before UNION SELECT * FROM after
	// "
	return 0, nil
}

func (q Query) UpdateRank(projectId, columnId rsrc.Id, rank rsrc.Rank) error {
	return nil
}

func (q Query) GetNextRank(projectId rsrc.Id, rank rsrc.Rank) (rsrc.Rank, error) {
	return "", nil
}

func (q Query) GetByName(projectId rsrc.Id, name string) (rsrc.Column, error) {
	return rsrc.Column{}, nil
}
