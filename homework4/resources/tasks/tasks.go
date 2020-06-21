package tasks

import (
	pg "github.com/AndreyKlimchuk/golang-learning/homework4/postgres"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type CreateRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
	rsrc.TaskSettableFields
}

type ReadRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
	TaskId    rsrc.Id
	expanded  bool
}

type ReadCollectionRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
}

type UpdateRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
	TaskId    rsrc.Id
	rsrc.TaskSettableFields
}

type DeleteRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
	TaskId    rsrc.Id
}

type UpdatePositionRequest struct {
	ProjectId   rsrc.Id
	ColumnId    rsrc.Id
	TaskId      rsrc.Id
	NewColumnId rsrc.Id
	AfterTaskId rsrc.Id
}

func (r CreateRequest) Create() (rsrc.Task, error) {
	if _, err := pg.Query().Columns().Get(r.ProjectId, r.ColumnId); err != nil {
		return rsrc.Task{}, rsrc.NewNotFoundOrInternalError("cannot get column", err)
	}
	tx, err := pg.Begin()
	if err != nil {
		return rsrc.Task{}, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer pg.Rollback(tx)
	maxRank, err := pg.QueryWithTX(tx).Tasks().GetAndBlockMaxRankByColumn(r.ColumnId)
	if err != nil {
		return rsrc.Task{}, rsrc.NewInternalError("cannot get max rank", err)
	}
	maxRank = rsrc.CalculateRank(maxRank, "")
	task, err := pg.QueryWithTX(tx).Tasks().Create(r.ProjectId, r.ColumnId, r.Name, r.Description, maxRank)
	if err != nil {
		return task, rsrc.NewInternalError("cannot create task", err)
	}
	if err := pg.Commit(tx); err != nil {
		return rsrc.Task{}, rsrc.NewInternalError("cannot commit transaction", err)
	}
	return task, nil
}

func (r ReadRequest) Read() (rsrc.Task, error) {
	return rsrc.Task{}, nil
}

func (r ReadCollectionRequest) ReadCollection() ([]rsrc.Task, error) {
	return []rsrc.Task{}, nil
}

func (r UpdateRequest) Update() error {
	return nil
}

func (r DeleteRequest) Delete() error {
	return nil
}

func (r UpdatePositionRequest) UpdatePosition() error {
	if _, err := pg.Query().Columns().Get(r.ProjectId, r.ColumnId); err != nil {
		return rsrc.NewNotFoundOrInternalError("cannot get column", err)
	}
	if r.ColumnId != r.NewColumnId {
		_, err := pg.Query().Columns().Get(r.ProjectId, r.NewColumnId)
		if pg.IsNoRowsError(err) {
			return rsrc.NewConflictError("column specified by new_column_id not found in project")
		} else {
			return rsrc.NewInternalError("cannot get new column", err)
		}
	}
	tx, err := pg.Begin()
	if err != nil {
		return rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer pg.Rollback(tx)
	// TODO
	if err := pg.Commit(tx); err != nil {
		return rsrc.NewInternalError("cannot commit transaction", err)
	}
	return nil
}
