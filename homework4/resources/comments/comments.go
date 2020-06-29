package comments

import (
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
)

type CreateRequest struct {
	TaskId common.Id
	common.CommentSettableFields
}

type ReadRequest struct {
	TaskId    common.Id
	CommentId common.Id
}

type ReadCollectionRequest struct {
	TaskId common.Id
}

type UpdateRequest struct {
	TaskId    common.Id
	CommentId common.Id
	common.CommentSettableFields
}

type DeleteRequest struct {
	TaskId    common.Id
	CommentId common.Id
}

func (r CreateRequest) Handle() (interface{}, error) {
	_, err := db.Query().Tasks().Get(r.TaskId)
	if err != nil {
		return common.Comment{}, common.NewNotFoundOrInternalError("cannot get task", err)
	}
	comment, err := db.Query().Comments().Create(r.TaskId, r.Text)
	return comment, common.MaybeNewInternalError("cannot create comment", err)
}

func (r ReadRequest) Handle() (interface{}, error) {
	comment, err := db.Query().Comments().Get(r.TaskId, r.CommentId)
	return comment, common.MaybeNewNotFoundOrInternalError("cannot get comment", err)
}

func (r ReadCollectionRequest) Handle() (interface{}, error) {
	comments, err := db.Query().Comments().GetMultiple(r.TaskId)
	return comments, common.MaybeNewInternalError("cannot read comments", err)
}

func (r UpdateRequest) Handle() (interface{}, error) {
	err := db.Query().Comments().Update(r.TaskId, r.CommentId, r.Text)
	return nil, common.MaybeNewNotFoundOrInternalError("cannot update comment", err)
}

func (r DeleteRequest) Handle() (interface{}, error) {
	err := db.Query().Comments().Delete(r.TaskId, r.CommentId)
	return nil, common.MaybeNewNotFoundOrInternalError("cannot delete comment", err)
}
