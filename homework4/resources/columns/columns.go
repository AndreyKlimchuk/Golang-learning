package columns

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type CreateRequest struct {
	ProjectId rsrc.Id
	rsrc.ColumnSettableFields
}

type ReadRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
}

type ReadCollectionRequest struct {
	ProjectId rsrc.Id
}

type UpdateRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
	rsrc.ColumnSettableFields
}

type UpdatePositionRequest struct {
	ProjectId     rsrc.Id `swaggerignore:"true"`
	ColumnId      rsrc.Id `swaggerignore:"true"`
	AfterColumnId rsrc.Id
}

type DeleteRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
}

func (r CreateRequest) Handle() (interface{}, error) {
	_, err := db.Query().Columns().GetByName(r.ProjectId, r.Name)
	if err == nil {
		return rsrc.Column{}, rsrc.NewConflictError("column with same name exists in project")
	} else if !common.IsNoRowsError(err) {
		return rsrc.Column{}, rsrc.NewInternalError("cannot get column by name", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return rsrc.Column{}, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	maxRank, err := db.QueryWithTX(tx).Columns().GetAndBlockMaxRank(r.ProjectId)
	if err != nil {
		return rsrc.Column{}, rsrc.NewNotFoundOrInternalError("cannot get max rank", err)
	}
	maxRank = rsrc.CalculateRank(maxRank, "")
	column, err := db.QueryWithTX(tx).Columns().Create(r.ProjectId, r.Name, maxRank)
	if err != nil {
		return rsrc.Column{}, rsrc.NewInternalError("cannot create column", err)
	}
	if err := db.Commit(tx); err != nil {
		return rsrc.Column{}, rsrc.NewInternalError("cannot commit transaction", err)
	}
	return column, nil
}

func (r ReadRequest) Handle() (interface{}, error) {
	column, err := db.Query().Columns().Get(r.ProjectId, r.ColumnId)
	return column, rsrc.MaybeNewNotFoundOrInternalError("cannot get column", err)
}

func (r ReadCollectionRequest) Handle() (interface{}, error) {
	columns, err := db.Query().Columns().GetMultiple(r.ProjectId)
	return columns, rsrc.MaybeNewInternalError("cannot get columns", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	column, err := db.Query().Columns().GetByName(r.ProjectId, r.Name)
	if err == nil && column.Id == r.ColumnId {
		return nil, rsrc.NewConflictError("column with specified name already exists in project")
	} else if err != nil && !common.IsNoRowsError(err) {
		return nil, rsrc.NewInternalError("cannot get column by name", err)
	}
	err = db.Query().Columns().Update(r.ProjectId, r.ColumnId, r.Name)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot update column", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	rank, err := db.QueryWithTX(tx).Columns().GetAndBlockRank(r.ProjectId, r.ColumnId)
	if err != nil {
		return nil, rsrc.NewNotFoundOrInternalError("cannot get column rank", err)
	}
	successorColumnId, err := db.QueryWithTX(tx).Columns().GetAndBlockSuccessorColumnId(r.ProjectId, rank)
	if common.IsNoRowsError(err) {
		return nil, rsrc.NewConflictError("project must contains at least one column")
	}
	if err != nil {
		return nil, rsrc.NewInternalError("cannot get successor column", err)
	}
	if err := moveTasks(tx, successorColumnId, r.ColumnId); err != nil {
		return nil, err
	}
	if err := db.QueryWithTX(tx).Columns().Delete(r.ColumnId); err != nil {
		return nil, rsrc.NewInternalError("cannot delete column", err)
	}
	if err := db.Commit(tx); err != nil {
		return nil, rsrc.NewInternalError("cannot commit transaction", err)
	}
	return nil, nil
}

func moveTasks(tx db.TX, dstColumnId rsrc.Id, srcColumnId rsrc.Id) error {
	tasksIds, err := db.QueryWithTX(tx).Tasks().GetAndBlockIdsByColumn(srcColumnId)
	if err != nil {
		return rsrc.NewInternalError("cannot get successor column tasks ids", err)
	}
	maxRank, err := db.QueryWithTX(tx).Tasks().GetAndBlockMaxRankByColumn(dstColumnId)
	if common.IsNoRowsError(err) {
		maxRank = ""
	} else if err != nil {
		return rsrc.NewInternalError("cannot get max task rank", err)
	}
	for _, taskId := range tasksIds {
		maxRank = rsrc.CalculateRank(maxRank, "")
		if err := db.QueryWithTX(tx).Tasks().UpdatePosition(taskId, dstColumnId, maxRank); err != nil {
			return rsrc.NewInternalError("cannot update task position", err)
		}
	}
	return nil
}

func (r UpdatePositionRequest) Handle() (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	var prevRank rsrc.Rank = ""
	if r.AfterColumnId > 0 {
		prevRank, err = db.QueryWithTX(tx).Columns().GetAndBlockRank(r.ProjectId, r.AfterColumnId)
		if common.IsNoRowsError(err) {
			return nil, rsrc.NewConflictError("column specified by after_column_id doesn't exists in project")
		}
		if err != nil {
			return nil, rsrc.NewInternalError("cannot get column specified by after_column_id", err)
		}
	}
	nextRank, err := db.QueryWithTX(tx).Columns().GetNextRank(r.ProjectId, prevRank)
	if common.IsNoRowsError(err) {
		nextRank = ""
	} else if err != nil {
		return nil, rsrc.NewInternalError("cannot get next column rank", err)
	}
	newRank := rsrc.CalculateRank(prevRank, nextRank)
	err = db.QueryWithTX(tx).Columns().UpdateRank(r.ProjectId, r.ColumnId, newRank)
	if err != nil {
		return nil, rsrc.NewNotFoundOrInternalError("cannot update column rank", err)
	}
	if err := db.Commit(tx); err != nil {
		return nil, rsrc.NewInternalError("cannot commit transaction", err)
	}
	return nil, nil
}
