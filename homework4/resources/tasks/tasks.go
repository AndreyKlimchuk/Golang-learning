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
	Expanded  bool
}

type UpdateRequest struct {
	TaskId rsrc.Id
	rsrc.TaskSettableFields
}

type DeleteRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
	TaskId    rsrc.Id
}

type UpdatePositionRequest struct {
	TaskId      rsrc.Id
	NewColumnId rsrc.Id
	AfterTaskId rsrc.Id
}

func (r CreateRequest) Create() (rsrc.Task, error) {
	if _, err := pg.Query().Columns().Get(r.ProjectId, r.ColumnId); err != nil {
		return rsrc.Task{}, rsrc.NewNotFoundOrInternalError("cannot get column", err)
	}
	tx, err := pg.Begin()
	defer pg.Rollback(tx)
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
	task, err := pg.Query().Tasks().Get(r.TaskId, r.Expanded)
	return task, rsrc.MaybeNewNotFoundOrInternalError("cannot read task", err)
}

func (r UpdateRequest) Update() error {
	err := pg.Query().Tasks().Update(r.TaskId, r.Name, r.Description)
	return rsrc.MaybeNewNotFoundOrInternalError("cannot update task", err)
}

func (r DeleteRequest) Delete() error {
	err := pg.Query().Tasks().Delete(r.TaskId)
	return rsrc.MaybeNewNotFoundOrInternalError("cannot delete task", err)
}

func (r UpdatePositionRequest) UpdatePosition() error {
	if err := validatePositionUpdate(r); err != nil {
		return err
	}
	tx, err := pg.Begin()
	if err != nil {
		return rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer pg.Rollback(tx)
	var prevRank rsrc.Rank = ""
	if r.AfterTaskId > 0 {
		prevRank, err = pg.QueryWithTX(tx).Tasks().GetAndBlockRank(r.NewColumnId, r.AfterTaskId)
		if pg.IsNoRowsError(err) {
			return rsrc.NewConflictError("task specified by after_task_id not found in target column")
		} else if err != nil {
			return rsrc.NewInternalError("cannot get previous task rank", err)
		}
	}
	nextRank, err := pg.QueryWithTX(tx).Tasks().GetNextRank(r.NewColumnId, prevRank)
	if pg.IsNoRowsError(err) {
		nextRank = ""
	} else if err != nil {
		return rsrc.NewInternalError("cannot get next task rank", err)
	}
	newRank := rsrc.CalculateRank(prevRank, nextRank)
	if err := pg.QueryWithTX(tx).Tasks().UpdatePosition(r.TaskId, r.NewColumnId, newRank); err != nil {
		return rsrc.NewNotFoundOrInternalError("cannot update task position", err)
	}
	if err := pg.Commit(tx); err != nil {
		return rsrc.NewInternalError("cannot commit transaction", err)
	}
	return nil
}

func validatePositionUpdate(r UpdatePositionRequest) error {
	task, err := pg.Query().Tasks().Get(r.TaskId, false)
	if err != nil {
		return rsrc.NewNotFoundOrInternalError("cannot get task", err)
	}
	if task.ColumnId != r.NewColumnId {
		_, err := pg.Query().Columns().Get(task.ProjectId, r.NewColumnId)
		if pg.IsNoRowsError(err) {
			return rsrc.NewConflictError("column specified by new_column_id not found in target project")
		} else if err != nil {
			return rsrc.NewInternalError("cannot get new column", err)
		}
	}
	return nil
}
