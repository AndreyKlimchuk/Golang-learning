package comments

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
)

type CreateRequest struct {
	TaskId rsrc.Id
	rsrc.CommentSettableFields
}

type ReadRequest struct {
	TaskId    rsrc.Id
	CommentId rsrc.Id
}

type ReadCollectionRequest struct {
	TaskId rsrc.Id
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

func (r CreateRequest) Handle() (interface{}, error) {
	_, err := db.Query().Tasks().Get(r.TaskId)
	if err != nil {
		return rsrc.Comment{}, rsrc.NewNotFoundOrInternalError("cannot get task", err)
	}
	comment, err := db.Query().Comments().Create(r.TaskId, r.Text)
	return comment, rsrc.MaybeNewInternalError("cannot create comment", err)
}

func (r ReadRequest) Handle() (interface{}, error) {
	comment, err := db.Query().Comments().Get(r.TaskId, r.CommentId)
	return comment, rsrc.MaybeNewNotFoundOrInternalError("cannot get comment", err)
}

func (r ReadCollectionRequest) Handle() (interface{}, error) {
	comments, err := db.Query().Comments().GetMultiple(r.TaskId)
	return comments, rsrc.MaybeNewInternalError("cannot read comments", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	err := db.Query().Comments().Update(r.TaskId, r.CommentId, r.Text)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot update comment", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	err := db.Query().Comments().Delete(r.TaskId, r.CommentId)
	return nil, rsrc.MaybeNewNotFoundOrInternalError("cannot delete comment", err)
}
