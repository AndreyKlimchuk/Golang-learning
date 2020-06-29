package comments

import (
	"context"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rcommon "github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
)

type QueryerWrap common.QueryerWrap

func (w QueryerWrap) Create(taskId rcommon.Id, text string) (rcommon.Comment, error) {
	comment := rcommon.Comment{CommentSettableFields: rcommon.CommentSettableFields{Text: text}}
	const q = "INSERT INTO comments (task_id, text, create_dt) VALUES ($1, $2, NOW()) RETURNING id"
	err := w.Q.QueryRow(context.Background(), q, taskId, text).Scan(&comment.Id)
	return comment, err
}

func (w QueryerWrap) Get(taskId, commentId rcommon.Id) (rcommon.Comment, error) {
	comment := rcommon.Comment{Id: commentId}
	const q = "SELECT text FROM comments WHERE task_id = $1 AND id = $2"
	err := w.Q.QueryRow(context.Background(), q, taskId, commentId).Scan(&comment.Text)
	return comment, err
}

func (w QueryerWrap) GetMultiple(taskId rcommon.Id) ([]rcommon.Comment, error) {
	comments := []rcommon.Comment{}
	const q = "SELECT id, text FROM comments WHERE task_id = $1 ORDER BY create_dt ASC"
	rows, err := w.Q.Query(context.Background(), q, taskId)
	if err != nil {
		return comments, err
	}
	defer rows.Close()
	c := rcommon.Comment{}
	for rows.Next() {
		err := rows.Scan(&c.Id, &c.Text)
		if err != nil {
			return comments, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (w QueryerWrap) Update(taskId, commentId rcommon.Id, text string) error {
	const q = "UPDATE comments SET text = $3 WHERE task_id = $1 AND id = $2"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, taskId, commentId, text))
}

func (w QueryerWrap) Delete(taskId, commentId rcommon.Id) error {
	const q = "DELETE FROM comments WHERE task_id = $1 AND id = $2"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, taskId, commentId))
}
