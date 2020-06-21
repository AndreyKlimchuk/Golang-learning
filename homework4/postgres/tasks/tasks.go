package tasks

import (
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type Query struct {
}

func (q Query) GetAndBlockIdsByColumn(columnId rsrc.Id) ([]rsrc.Id, error) {
	return []rsrc.Id{}, nil
}

func (q Query) GetAndBlockMaxRankByColumn(columnId rsrc.Id) (rsrc.Rank, error) {
	return "", nil
}

func (q Query) UpdatePosition(taskId, columnId rsrc.Id, rank rsrc.Rank) error {
	return nil
}

func (q Query) Get(taskId rsrc.Id) (rsrc.Task, error) {
	return rsrc.Task{}, nil
}

func (q Query) Create(projectId, columnId rsrc.Id, name string, description string, rank rsrc.Rank) (rsrc.Task, error) {
	return rsrc.Task{}, nil
}
