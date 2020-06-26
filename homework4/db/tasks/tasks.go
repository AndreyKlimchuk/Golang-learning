package tasks

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type QueryerWrap struct {
	Q common.Queryer
}

func (w QueryerWrap) GetAndBlockIdsByColumn(columnId rsrc.Id) ([]rsrc.Id, error) {
	return []rsrc.Id{}, nil
}

func (w QueryerWrap) GetAndBlockMaxRankByColumn(columnId rsrc.Id) (rsrc.Rank, error) {
	return "", nil
}

func (w QueryerWrap) UpdatePosition(taskId, columnId rsrc.Id, rank rsrc.Rank) error {
	return nil
}

func (w QueryerWrap) Get(taskId rsrc.Id) (rsrc.Task, error) {
	return rsrc.Task{}, nil
}

func (w QueryerWrap) GetExpanded(taskId rsrc.Id) (rsrc.TaskExpanded, error) {
	return rsrc.TaskExpanded{}, nil
}

func (w QueryerWrap) Create(projectId, columnId rsrc.Id, name string, description string, rank rsrc.Rank) (rsrc.Task, error) {
	return rsrc.Task{}, nil
}

func (w QueryerWrap) GetAndBlockRank(columnId, taskId rsrc.Id) (rsrc.Rank, error) {
	return "", nil
}

func (w QueryerWrap) GetNextRank(columnId rsrc.Id, rank rsrc.Rank) (rsrc.Rank, error) {
	return "", nil
}

func (w QueryerWrap) Update(taskId rsrc.Id, name string, description string) error {
	return nil
}

func (w QueryerWrap) Delete(taskId rsrc.Id) error {
	return nil
}
