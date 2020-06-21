package comments

import (
	pg "github.com/AndreyKlimchuk/golang-learning/homework4/postgres"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type CreateRequest struct {
	TaskId    rsrc.Id
	rsrc.CommentSettableFields
}

type ReadRequest struct {
	TaskId    rsrc.Id
	CommentId rsrc.Id
}

type ReadCollectionRequest struct {
	TaskId    rsrc.Id
}

type UpdateRequest struct {
	TaskId    rsrc.Id
	CommentId rsrc.Id
	rsrc.CommentSettableFields
}

type DeleteRequest struct {
	ProjectId rsrc.Id
	ColumnId  rsrc.Id
	TaskId    rsrc.Id
	CommentId rsrc.Id
}

func (r CreateRequest) Create() (rsrc.Comment, error) {
	_, err := pg.Query().Tasks().Get(r.TaskId)
	if err != nil {
		return rsrc.Comment{}, rsrc.NewNotFoundOrInternalError("cannot get task", err)
	}
	comment, err := pg.Query().Comments().Create(r.TaskId, r.Text)
	return comment, rsrc.MaybeNewInternalError("cannot create comment", err)
}

func (r ReadRequest) Read() (rsrc.Comment, error) {
	comment, err := pg.Query().Comments().Read(r.TaskId, r.CommentId)
	return comment, rsrc.MaybeNewNotFoundOrInternalError("cannot get comment", err)
}

func (r ReadCollectionRequest) ReadCollection() ([]rsrc.Comment, error) {
	comments, err := pg.Query().Comments().ReadMultiple(r.TaskId)
	return comments, rsrc.MaybeNewInternalError("cannot read comments", err)
}

func (r UpdateRequest) Update() error {
	err := pg.Query().Comments().Update(r.TaskId, r.CommentId, r.Text)
	return rsrc.MaybeNewNotFoundOrInternalError("cannot update comment", err)
}

func (r DeleteRequest) Delete() error {
	err := pg.Query().Comments().Delete(r.TaskId, r.CommentId)
	return rsrc.MaybeNewNotFoundOrInternalError("cannot delete comment", err)
}
