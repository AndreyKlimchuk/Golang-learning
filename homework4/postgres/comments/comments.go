package comments

import rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"

type Query struct {
}

func (q Query) Create(taskId rsrc.Id, text string) (rsrc.Comment, error) {
	return rsrc.Comment{}, nil
}

func (q Query) Get(taskId, commentId rsrc.Id) (rsrc.Comment, error) {
	return rsrc.Comment{}, nil
}

func (q Query) GetMultiple(taskId rsrc.Id) ([]rsrc.Comment, error) {
	return []rsrc.Comment{}, nil
}

func (q Query) Update(taskId, commentId rsrc.Id, text string) error {
	return nil
}

func (q Query) Delete(taskId, commentId rsrc.Id) error {
	return nil
}
