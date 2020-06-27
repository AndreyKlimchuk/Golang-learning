package tasks

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type CreateRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
	rsrc.TaskSettableFields
}

type ReadRequest struct {
	TaskId   rsrc.Id
	Expanded bool
}

type UpdateRequest struct {
	TaskId rsrc.Id
	rsrc.TaskSettableFields
}

type DeleteRequest struct {
	TaskId rsrc.Id
}

type UpdatePositionRequest struct {
	TaskId      rsrc.Id `swaggerignore:"true"`
	NewColumnId rsrc.Id `json:"new_column_id"`
	AfterTaskId rsrc.Id `json:"after_task_id"`
}

func (r CreateRequest) Handle() (interface{}, error) {
	if _, err := db.Query().Columns().Get(r.ProjectId, r.ColumnId); err != nil {
		return rsrc.Task{}, rsrc.NewNotFoundOrInternalError("cannot get column", err)
	}
	tx, err := db.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return rsrc.Task{}, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	maxRank, err := db.QueryWithTX(tx).Tasks().GetAndBlockMaxRankByColumn(r.ColumnId)
	if err != nil {
		return rsrc.Task{}, rsrc.NewInternalError("cannot get max rank", err)
	}
	maxRank = rsrc.CalculateRank(maxRank, "")
	task, err := db.QueryWithTX(tx).Tasks().Create(r.ProjectId, r.ColumnId, r.Name, r.Description, maxRank)
	if err != nil {
		return task, rsrc.NewInternalError("cannot create task", err)
	}
	if err := db.Commit(tx); err != nil {
		return rsrc.Task{}, rsrc.NewInternalError("cannot commit transaction", err)
	}
	return task, nil
}

func (r ReadRequest) Handle() (interface{}, error) {
	var task rsrc.Resource
	var err error
	if r.Expanded {
		task, err = db.Query().Tasks().GetExpanded(r.TaskId)
	} else {
		task, err = db.Query().Tasks().Get(r.TaskId)
	}
	return task, rsrc.MaybeNewNotFoundOrInternalError("cannot read task", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	err := db.Query().Tasks().Update(r.TaskId, r.Name, r.Description)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot update task", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	err := db.Query().Tasks().Delete(r.TaskId)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot delete task", err)
}

func (r UpdatePositionRequest) Handle() (interface{}, error) {
	if err := validatePositionUpdate(r); err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	var prevRank rsrc.Rank = ""
	if r.AfterTaskId > 0 {
		prevRank, err = db.QueryWithTX(tx).Tasks().GetAndBlockRank(r.NewColumnId, r.AfterTaskId)
		if common.IsNoRowsError(err) {
			return nil, rsrc.NewConflictError("task specified by after_task_id not found in target column")
		} else if err != nil {
			return nil, rsrc.NewInternalError("cannot get previous task rank", err)
		}
	}
	nextRank, err := db.QueryWithTX(tx).Tasks().GetNextRank(r.NewColumnId, prevRank)
	if common.IsNoRowsError(err) {
		nextRank = ""
	} else if err != nil {
		return nil, rsrc.NewInternalError("cannot get next task rank", err)
	}
	newRank := rsrc.CalculateRank(prevRank, nextRank)
	if err := db.QueryWithTX(tx).Tasks().UpdatePosition(r.TaskId, r.NewColumnId, newRank); err != nil {
		return nil, rsrc.NewNotFoundOrInternalError("cannot update task position", err)
	}
	if err := db.Commit(tx); err != nil {
		return nil, rsrc.NewInternalError("cannot commit transaction", err)
	}
	return nil, nil
}

func validatePositionUpdate(r UpdatePositionRequest) error {
	task, err := db.Query().Tasks().Get(r.TaskId)
	if err != nil {
		return rsrc.NewNotFoundOrInternalError("cannot get task", err)
	}
	if task.ColumnId != r.NewColumnId {
		_, err := db.Query().Columns().Get(task.ProjectId, r.NewColumnId)
		if common.IsNoRowsError(err) {
			return rsrc.NewConflictError("column specified by new_column_id not found in target project")
		} else if err != nil {
			return rsrc.NewInternalError("cannot get new column", err)
		}
	}
	return nil
}
