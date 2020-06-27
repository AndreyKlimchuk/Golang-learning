package projects

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

const defaultColumnName string = "default"

type CreateRequest struct {
	rsrc.ProjectSettableFields
}

type ReadRequest struct {
	ProjectId rsrc.Id
	Expanded  bool
}

type ReadCollectionRequest struct {
}

type UpdateRequest struct {
	ProjectId rsrc.Id
	rsrc.ProjectSettableFields
}

type DeleteRequest struct {
	ProjectId rsrc.Id
}

func (r CreateRequest) Handle() (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return rsrc.Project{}, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer db.Rollback(tx)
	project, err := db.QueryWithTX(tx).Projects().Create(r.Name, r.Description)
	if err != nil {
		return rsrc.Project{}, rsrc.NewInternalError("cannot create project", err)
	}
	rank := rsrc.CalculateRank("", "")
	column, err := db.QueryWithTX(tx).Columns().Create(project.Id, defaultColumnName, rank)
	if err != nil {
		return rsrc.Project{}, rsrc.NewInternalError("cannot create column", err)
	}
	if err := db.Commit(tx); err != nil {
		return rsrc.Project{}, rsrc.NewInternalError("cannot commit transaction", err)
	}
	projectExpanded := rsrc.ProjectExpanded{Project: project, Columns: []rsrc.ColumnExpanded{column}}
	return projectExpanded, nil
}

func (r ReadRequest) Handle() (interface{}, error) {
	var project rsrc.Resource
	var err error
	if r.Expanded {
		project, err = db.Query().Projects().GetExpanded(r.ProjectId)
	} else {
		project, err = db.Query().Projects().Get(r.ProjectId)
	}
	return project, rsrc.MaybeNewNotFoundOrInternalError("cannot get project", err)
}

func (_ ReadCollectionRequest) Handle() (interface{}, error) {
	project, err := db.Query().Projects().GetMultiple()
	return project, rsrc.MaybeNewInternalError("cannot get projects", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	err := db.Query().Projects().Update(r.ProjectId, r.Name, r.Description)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot update project", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	err := db.Query().Projects().Delete(r.ProjectId)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot delete project", err)
}
