package projects

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
)

type CreateRequest struct {
	common.ProjectSettableFields
}

type ReadRequest struct {
	ProjectId common.Id
	Expanded  bool
}

type ReadCollectionRequest struct {
}

type UpdateRequest struct {
	ProjectId common.Id
	common.ProjectSettableFields
}

type DeleteRequest struct {
	ProjectId common.Id
}

func (r CreateRequest) Handle() (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return common.Project{}, common.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	project, err := db.QueryWithTX(tx).Projects().Create(r.Name, r.Description)
	if err != nil {
		return common.Project{}, common.NewInternalError("cannot create project", err)
	}
	rank := common.CalculateRankInitial()
	column, err := db.QueryWithTX(tx).Columns().Create(project.Id, common.DefaultColumnName, rank)
	if err != nil {
		return common.Project{}, common.NewInternalError("cannot create column", err)
	}
	if err := db.Commit(tx); err != nil {
		return common.Project{}, common.NewInternalError("cannot commit transaction", err)
	}
	projectExpanded := common.ProjectExpanded{Project: project, Columns: []common.ColumnExpanded{column}}
	return projectExpanded, nil
}

func (r ReadRequest) Handle() (interface{}, error) {
	var project resources.Resource
	var err error
	if r.Expanded {
		project, err = db.Query().Projects().GetExpanded(r.ProjectId)
	} else {
		project, err = db.Query().Projects().Get(r.ProjectId)
	}
	return project, common.MaybeNewNotFoundOrInternalError("cannot get project", err)
}

func (_ ReadCollectionRequest) Handle() (interface{}, error) {
	project, err := db.Query().Projects().GetMultiple()
	return project, common.MaybeNewInternalError("cannot get projects", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	err := db.Query().Projects().Update(r.ProjectId, r.Name, r.Description)
	return nil, common.MaybeNewNotFoundOrInternalError("cannot update project", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	err := db.Query().Projects().Delete(r.ProjectId)
	return nil, common.MaybeNewNotFoundOrInternalError("cannot delete project", err)
}
