package projects

import (
	pg "github.com/AndreyKlimchuk/golang-learning/homework4/postgres"
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
	tx, err := pg.Begin()
	if err != nil {
		return rsrc.Project{}, rsrc.NewInternalError("cannot begin transaction", err)
	}
	defer pg.Rollback(tx)
	project, err := pg.QueryWithTX(tx).Projects().Create(r.Name, r.Description)
	if err != nil {
		return rsrc.Project{}, rsrc.NewInternalError("cannot create project", err)
	}
	rank := rsrc.CalculateRank("", "")
	column, err := pg.QueryWithTX(tx).Columns().Create(project.Id, defaultColumnName, rank)
	if err != nil {
		return rsrc.Project{}, rsrc.NewInternalError("cannot create column", err)
	}
	project.Columns = []rsrc.Column{column}
	if err := pg.Commit(tx); err != nil {
		return rsrc.Project{}, rsrc.NewInternalError("cannot commit transaction", err)
	}
	return project, nil
}

func (r ReadRequest) Handle() (interface{}, error) {
	project, err := pg.Query().Projects().Get(r.ProjectId, r.Expanded)
	return project, rsrc.MaybeNewNotFoundOrInternalError("cannot get project", err)
}

func (_ ReadCollectionRequest) Handle() (interface{}, error) {
	project, err := pg.Query().Projects().GetMultiple()
	return project, rsrc.MaybeNewInternalError("cannot get projects", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	err := pg.Query().Projects().Update(r.ProjectId, r.Name, r.Description)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot update project", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	err := pg.Query().Projects().Delete(r.ProjectId)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot delete project", err)
}
