package columns

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rcommon "github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
)

type CreateRequest struct {
	ProjectId rcommon.Id
	rcommon.ColumnSettableFields
}

type ReadRequest struct {
	ProjectId rcommon.Id
	ColumnId  rcommon.Id
}

type ReadCollectionRequest struct {
	ProjectId rcommon.Id
}

type UpdateRequest struct {
	ProjectId rcommon.Id
	ColumnId  rcommon.Id
	rcommon.ColumnSettableFields
}

type UpdatePositionRequest struct {
	ProjectId rcommon.Id
	ColumnId  rcommon.Id `validate:"nefield=UpdatePositionRequestBody.AfterColumnId"`
	UpdatePositionRequestBody
}

type UpdatePositionRequestBody struct {
	AfterColumnId rcommon.Id `swaggertype:"primitive,integer"`
}

type DeleteRequest struct {
	ProjectId rcommon.Id
	ColumnId  rcommon.Id
}

func (r CreateRequest) Handle() (interface{}, error) {
	_, err := db.Query().Columns().GetByName(r.ProjectId, r.Name)
	if err == nil {
		return rcommon.Column{}, rcommon.NewConflictError("column with same name exists in project")
	} else if !common.IsNoRowsError(err) {
		return rcommon.Column{}, rcommon.NewInternalError("cannot get column by name", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return rcommon.Column{}, rcommon.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	maxRank, err := db.QueryWithTX(tx).Columns().GetAndBlockMaxRank(r.ProjectId)
	if err != nil {
		return rcommon.Column{}, rcommon.NewNotFoundOrInternalError("cannot get max rank", err)
	}
	maxRank = rcommon.CalculateRankHigher(maxRank)
	column, err := db.QueryWithTX(tx).Columns().Create(r.ProjectId, r.Name, maxRank)
	if err != nil {
		return rcommon.Column{}, rcommon.NewInternalError("cannot create column", err)
	}
	if err := db.Commit(tx); err != nil {
		return rcommon.Column{}, rcommon.NewInternalError("cannot commit transaction", err)
	}
	return column, nil
}

func (r ReadRequest) Handle() (interface{}, error) {
	column, err := db.Query().Columns().Get(r.ProjectId, r.ColumnId)
	return column, rcommon.MaybeNewNotFoundOrInternalError("cannot get column", err)
}

func (r ReadCollectionRequest) Handle() (interface{}, error) {
	columns, err := db.Query().Columns().GetMultiple(r.ProjectId)
	return columns, rcommon.MaybeNewInternalError("cannot get columns", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	_, err := db.Query().Columns().GetByName(r.ProjectId, r.Name)
	if err == nil {
		return nil, rcommon.NewConflictError("column with specified name already exists in project")
	} else if !common.IsNoRowsError(err) {
		return nil, rcommon.NewInternalError("cannot get column by name", err)
	}
	err = db.Query().Columns().Update(r.ProjectId, r.ColumnId, r.Name)
	return nil, rcommon.MaybeNewNotFoundOrInternalError("cannot update column", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, rcommon.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	rank, err := db.QueryWithTX(tx).Columns().GetAndBlockRank(r.ProjectId, r.ColumnId)
	if err != nil {
		return nil, rcommon.NewNotFoundOrInternalError("cannot get column rank", err)
	}
	successorColumnId, err := db.QueryWithTX(tx).Columns().GetAndBlockSuccessorColumnId(r.ProjectId, rank)
	if common.IsNoRowsError(err) {
		return nil, rcommon.NewConflictError("project must contains at least one column")
	}
	if err != nil {
		return nil, rcommon.NewInternalError("cannot get successor column", err)
	}
	if err := moveTasks(tx, successorColumnId, r.ColumnId); err != nil {
		return nil, err
	}
	if err := db.QueryWithTX(tx).Columns().Delete(r.ColumnId); err != nil {
		return nil, rcommon.NewInternalError("cannot delete column", err)
	}
	if err := db.Commit(tx); err != nil {
		return nil, rcommon.NewInternalError("cannot commit transaction", err)
	}
	return nil, nil
}

func moveTasks(tx db.TX, dstColumnId rcommon.Id, srcColumnId rcommon.Id) error {
	tasksIds, err := db.QueryWithTX(tx).Tasks().GetAndBlockIdsByColumn(srcColumnId)
	if err != nil {
		return rcommon.NewInternalError("cannot get successor column tasks ids", err)
	}
	maxRank, err := db.QueryWithTX(tx).Tasks().GetAndBlockMaxRankByColumn(dstColumnId)
	if common.IsNoRowsError(err) {
		maxRank = ""
	} else if err != nil {
		return rcommon.NewInternalError("cannot get max task rank", err)
	}
	for _, taskId := range tasksIds {
		maxRank = rcommon.CalculateRankHigher(maxRank)
		if err := db.QueryWithTX(tx).Tasks().UpdatePosition(taskId, dstColumnId, maxRank); err != nil {
			return rcommon.NewInternalError("cannot update task position", err)
		}
	}
	return nil
}

func (r UpdatePositionRequest) Handle() (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, rcommon.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	var prevRank rcommon.Rank = ""
	if r.AfterColumnId > 0 {
		prevRank, err = db.QueryWithTX(tx).Columns().GetAndBlockRank(r.ProjectId, r.AfterColumnId)
		if common.IsNoRowsError(err) {
			return nil, rcommon.NewConflictError("column specified by after_column_id doesn't exists in project")
		}
		if err != nil {
			return nil, rcommon.NewInternalError("cannot get column specified by after_column_id", err)
		}
	}
	var newRank rcommon.Rank
	nextRank, err := db.QueryWithTX(tx).Columns().GetNextRank(r.ProjectId, prevRank)
	if err == nil {
		newRank = rcommon.CalculateRankBetween(prevRank, nextRank)
	} else if common.IsNoRowsError(err) {
		newRank = rcommon.CalculateRankHigher(prevRank)
	} else if err != nil {
		return nil, rcommon.NewInternalError("cannot get next column rank", err)
	}
	err = db.QueryWithTX(tx).Columns().UpdateRank(r.ProjectId, r.ColumnId, newRank)
	if err != nil {
		return nil, rcommon.NewNotFoundOrInternalError("cannot update column rank", err)
	}
	if err := db.Commit(tx); err != nil {
		return nil, rcommon.NewInternalError("cannot commit transaction", err)
	}
	return nil, nil
}
