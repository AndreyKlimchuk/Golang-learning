package projects

import (
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

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

func (r CreateRequest) Create() (rsrc.Project, error) {
	return rsrc.Project{}, nil
}

func (r ReadRequest) Read() (rsrc.Project, error) {
	return rsrc.Project{}, nil
}

func (r ReadCollectionRequest) ReadCollection() ([]rsrc.Project, error) {
	return []rsrc.Project{}, nil
}

func (r UpdateRequest) Update() error {
	return nil
}

func (r DeleteRequest) Delete() error {
	return nil
}
