package tasks

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources"
	rcommon "github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
)

type CreateRequest struct {
	ProjectId rcommon.Id
	ColumnId  rcommon.Id
	rcommon.TaskSettableFields
}

type ReadRequest struct {
	TaskId   rcommon.Id
	Expanded bool
}

type UpdateRequest struct {
	TaskId rcommon.Id
	rcommon.TaskSettableFields
}

type DeleteRequest struct {
	TaskId rcommon.Id
}

type UpdatePositionRequest struct {
	TaskId rcommon.Id `validate:"nefield=UpdatePositionRequestBody.AfterTaskId"`
	UpdatePositionRequestBody
}

type UpdatePositionRequestBody struct {
	NewColumnId rcommon.Id `json:"new_column_id" swaggertype:"primitive,integer"`
	AfterTaskId rcommon.Id `json:"after_task_id" swaggertype:"primitive,integer"`
}

func (r CreateRequest) Handle() (interface{}, error) {
	if _, err := db.Query().Columns().Get(r.ProjectId, r.ColumnId); err != nil {
		return rcommon.Task{}, rcommon.NewNotFoundOrInternalError("cannot get column", err)
	}
	tx, err := db.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return rcommon.Task{}, rcommon.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	maxRank, err := db.QueryWithTX(tx).Tasks().GetAndBlockMaxRankByColumn(r.ColumnId)
	if common.IsNoRowsError(err) {
		maxRank = ""
	} else if err != nil {
		return rcommon.Task{}, rcommon.NewInternalError("cannot get max rank", err)
	}
	maxRank = rcommon.CalculateRankHigher(maxRank)
	task, err := db.QueryWithTX(tx).Tasks().Create(r.ProjectId, r.ColumnId, r.Name, r.Description, maxRank)
	if err != nil {
		return task, rcommon.NewInternalError("cannot create task", err)
	}
	if err := db.Commit(tx); err != nil {
		return rcommon.Task{}, rcommon.NewInternalError("cannot commit transaction", err)
	}
	return task, nil
}

func (r ReadRequest) Handle() (interface{}, error) {
	var task resources.Resource
	var err error
	if r.Expanded {
		task, err = db.Query().Tasks().GetExpanded(r.TaskId)
	} else {
		task, err = db.Query().Tasks().Get(r.TaskId)
	}
	return task, rcommon.MaybeNewNotFoundOrInternalError("cannot read task", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	err := db.Query().Tasks().Update(r.TaskId, r.Name, r.Description)
	return nil, rcommon.MaybeNewNotFoundOrInternalError("cannot update task", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	err := db.Query().Tasks().Delete(r.TaskId)
	return nil, rcommon.MaybeNewNotFoundOrInternalError("cannot delete task", err)
}

func (r UpdatePositionRequest) Handle() (interface{}, error) {
	if err := validatePositionUpdate(r); err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, rcommon.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	var prevRank rcommon.Rank = ""
	if r.AfterTaskId > 0 {
		prevRank, err = db.QueryWithTX(tx).Tasks().GetAndBlockRank(r.NewColumnId, r.AfterTaskId)
		if common.IsNoRowsError(err) {
			return nil, rcommon.NewConflictError("task specified by after_task_id not found in target column")
		} else if err != nil {
			return nil, rcommon.NewInternalError("cannot get previous task rank", err)
		}
	}
	var newRank rcommon.Rank
	nextRank, err := db.QueryWithTX(tx).Tasks().GetNextRank(r.NewColumnId, prevRank)
	if err == nil {
		newRank = rcommon.CalculateRankBetween(prevRank, nextRank)
	} else if common.IsNoRowsError(err) {
		newRank = rcommon.CalculateRankHigher(prevRank)
	} else if err != nil {
		return nil, rcommon.NewInternalError("cannot get next task rank", err)
	}
	if err := db.QueryWithTX(tx).Tasks().UpdatePosition(r.TaskId, r.NewColumnId, newRank); err != nil {
		return nil, rcommon.NewNotFoundOrInternalError("cannot update task position", err)
	}
	if err := db.Commit(tx); err != nil {
		return nil, rcommon.NewInternalError("cannot commit transaction", err)
	}
	return nil, nil
}

func validatePositionUpdate(r UpdatePositionRequest) error {
	task, err := db.Query().Tasks().Get(r.TaskId)
	if err != nil {
		return rcommon.NewNotFoundOrInternalError("cannot get task", err)
	}
	if task.ColumnId != r.NewColumnId {
		_, err := db.Query().Columns().Get(task.ProjectId, r.NewColumnId)
		if common.IsNoRowsError(err) {
			return rcommon.NewConflictError("column specified by new_column_id not found in target project")
		} else if err != nil {
			return rcommon.NewInternalError("cannot get new column", err)
		}
	}
	return nil
}
