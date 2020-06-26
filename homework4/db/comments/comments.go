package comments

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type QueryerWrap struct {
	Q common.Queryer
}

func (w QueryerWrap) Create(taskId rsrc.Id, text string) (rsrc.Comment, error) {
	return rsrc.Comment{}, nil
}

func (w QueryerWrap) Get(taskId, commentId rsrc.Id) (rsrc.Comment, error) {
	return rsrc.Comment{}, nil
}

func (w QueryerWrap) GetMultiple(taskId rsrc.Id) ([]rsrc.Comment, error) {
	return []rsrc.Comment{}, nil
}

func (w QueryerWrap) Update(taskId, commentId rsrc.Id, text string) error {
	return nil
}

func (w QueryerWrap) Delete(taskId, commentId rsrc.Id) error {
	return nil
}
