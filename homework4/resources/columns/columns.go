package comments

import (
	pg "github.com/AndreyKlimchuk/golang-learning/homework4/postgres"
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
	ProjectId     rsrc.Id
	ColumnId      rsrc.Id
	AfterColumnId rsrc.Id
}

type DeleteRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
}

func (r CreateRequest) Create() (rsrc.Column, error) {
	column, err := pg.Query().Columns().GetByName(r.ProjectId, r.Name)
	if err == nil {
		return rsrc.Column{}, rsrc.NewConflictError("column with same name exists in project")
	} else if !pg.IsNoRowsError(err) {
		return rsrc.Column{}, rsrc.NewInternalError("cannot get column by name", err)
	}
	tx, err := pg.Begin()
	if err != nil {
		return rsrc.Column{}, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer pg.Rollback(tx)
	maxRank, err := pg.QueryWithTX(tx).Columns().GetAndBlockMaxRank(r.ProjectId)
	if err != nil {
		return rsrc.Column{}, rsrc.NewNotFoundOrInternalError("cannot get max rank", err)
	}
	maxRank = rsrc.CalculateRank(maxRank, "")
	column, err = pg.QueryWithTX(tx).Columns().Create(r.ProjectId, r.Name, maxRank)
	if err != nil {
		return rsrc.Column{}, rsrc.NewInternalError("cannot create column", err)
	}
	if err := pg.Commit(tx); err != nil {
		return rsrc.Column{}, rsrc.NewInternalError("cannot commit transaction", err)
	}
	return column, nil
}

func (r ReadRequest) Read() (rsrc.Column, error) {
	column, err := pg.Query().Columns().Get(r.ProjectId, r.ColumnId)
	return column, rsrc.MaybeNewNotFoundOrInternalError("cannot get column", err)
}

func (r ReadCollectionRequest) ReadCollection() ([]rsrc.Column, error) {
	columns, err := pg.Query().Columns().GetMultiple(r.ProjectId)
	return columns, rsrc.MaybeNewInternalError("cannot get columns", err)
}

func (r UpdateRequest) Update() error {
	column, err := pg.Query().Columns().GetByName(r.ProjectId, r.Name)
	if err == nil && column.Id == r.ColumnId {
		return rsrc.NewConflictError("column with specified name already exists in project")
	} else if err != nil && !pg.IsNoRowsError(err) {
		return rsrc.NewInternalError("cannot get column by name", err)
	}
	err = pg.Query().Columns().Update(r.ProjectId, r.ColumnId, r.Name)
	return rsrc.MaybeNewNotFoundOrInternalError("cannot update column", err)
}

func (r DeleteRequest) Delete() error {
	tx, err := pg.Begin()
	if err != nil {
		return rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer pg.Rollback(tx)
	rank, err := pg.QueryWithTX(tx).Columns().GetAndBlockRank(r.ProjectId, r.ColumnId)
	if err != nil {
		return rsrc.NewNotFoundOrInternalError("cannot get column rank", err)
	}
	successorColumnId, err := pg.QueryWithTX(tx).Columns().GetAndBlockSuccessorColumnId(r.ProjectId, rank)
	if pg.IsNoRowsError(err)  {
		return rsrc.NewConflictError("project must contains at least one column")
	}
	if err != nil {
		return rsrc.NewInternalError("cannot get successor column", err)
	}
	if err := moveTasks(tx, successorColumnId, r.ColumnId); err != nil {
		return err
	}
	if err := pg.QueryWithTX(tx).Columns().Delete(r.ColumnId); err != nil {
		return rsrc.NewInternalError("cannot delete column", err)
	}
	if err := pg.Commit(tx); err != nil {
		return rsrc.NewInternalError("cannot commit transaction", err)
	}
	return nil
}

func moveTasks(tx pg.TX, dstColumnId rsrc.Id, srcColumnId rsrc.Id) error {
	tasksIds, err := pg.QueryWithTX(tx).Tasks().GetAndBlockIdsByColumn(srcColumnId)
	if err != nil {
		return rsrc.NewInternalError("cannot get successor column tasks ids", err)
	}
	maxRank, err := pg.QueryWithTX(tx).Tasks().GetAndBlockMaxRankByColumn(dstColumnId)
	if pg.IsNoRowsError(err) {
		maxRank = ""
	} else if err != nil {
		return rsrc.NewInternalError("cannot get max task rank", err)
	}
	for _, taskId := range tasksIds {
		maxRank = rsrc.CalculateRank(maxRank, "")
		if err := pg.QueryWithTX(tx).Tasks().UpdatePosition(taskId, dstColumnId, maxRank); err != nil {
			return rsrc.NewInternalError("cannot update task position", err)
		}
	}
	return nil
}

func (r UpdatePositionRequest) UpdatePosition() error {
	tx, err := pg.Begin()
	if err != nil {
		return rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer pg.Rollback(tx)
	var prevRank rsrc.Rank = ""
	if r.AfterColumnId > 0 {
		prevRank, err = pg.QueryWithTX(tx).Columns().GetAndBlockRank(r.ProjectId, r.AfterColumnId)
		if pg.IsNoRowsError(err) {
			return rsrc.NewConflictError("column specified by after_column_id doesn't exists in project")
		}
		if err != nil {
			return rsrc.NewInternalError("cannot get column specified by after_column_id", err)
		}
	}
	nextRank, err := pg.QueryWithTX(tx).Columns().GetNextRank(r.ProjectId, prevRank)
	if pg.IsNoRowsError(err) {
		nextRank = ""
	} else if err != nil {
		return rsrc.NewInternalError("cannot get next column rank", err)
	}
	newRank := rsrc.CalculateRank(prevRank, nextRank)
	err = pg.QueryWithTX(tx).Columns().UpdateRank(r.ProjectId, r.ColumnId, newRank)
	if err != nil {
		return rsrc.NewNotFoundOrInternalError("cannot update column rank", err)
	}
	if err := pg.Commit(tx); err != nil {
		return rsrc.NewInternalError("cannot commit transaction", err)
	}
	return nil
}
