package projects

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type QueryerWrap struct {
	Q common.Queryer
}

func (w QueryerWrap) Create(name string, description string) (rsrc.ProjectExpanded, error) {
	return rsrc.ProjectExpanded{}, nil
}

func (w QueryerWrap) Get(projectId rsrc.Id) (rsrc.Project, error) {
	return rsrc.Project{}, nil
}

func (w QueryerWrap) GetExpanded(projectId rsrc.Id) (rsrc.ProjectExpanded, error) {
	return rsrc.ProjectExpanded{}, nil
}

func (w QueryerWrap) GetMultiple() ([]rsrc.Project, error) {
	return []rsrc.Project{}, nil
}

func (w QueryerWrap) Update(projectId rsrc.Id, name string, description string) error {
	return nil
}

func (w QueryerWrap) Delete(projectId rsrc.Id) error {
	return nil
}
