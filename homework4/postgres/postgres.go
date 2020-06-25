package postgres

import (
	"database/sql"
	"errors"

	"github.com/AndreyKlimchuk/golang-learning/homework4/postgres/columns"
	"github.com/AndreyKlimchuk/golang-learning/homework4/postgres/comments"
	"github.com/AndreyKlimchuk/golang-learning/homework4/postgres/projects"
	"github.com/AndreyKlimchuk/golang-learning/homework4/postgres/tasks"
)

type TX struct {
}

type query struct {
}

func (q query) Projects() projects.Query {
	return projects.Query{}
}

func (q query) Columns() columns.Query {
	return columns.Query{}
}

func (q query) Tasks() tasks.Query {
	return tasks.Query{}
}

func (q query) Comments() comments.Query {
	return comments.Query{}
}

func QueryWithTX(tx TX) query {
	return query{}
}

func Query() query {
	return query{}
}

func Begin() (TX, error) {
	return TX{}, nil
}

func Commit(tx TX) error {
	return nil
}

func Rollback(tx TX) {

}

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
