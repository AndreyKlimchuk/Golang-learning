package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AndreyKlimchuk/golang-learning/homework4/logger"
	"go.uber.org/zap"

	"github.com/AndreyKlimchuk/golang-learning/homework4/db/columns"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/comments"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/projects"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/tasks"
	"github.com/jackc/pgx/v4/pgxpool"
)

type queryerWrap struct {
	q common.Queryer
}

type TX pgx.Tx

var pool *pgxpool.Pool

func Init(conn string) (err error) {
	pool, err = pgxpool.Connect(context.Background(), conn)
	if err != nil {
		return err
	}
	return nil
}

func (w queryerWrap) Projects() projects.QueryerWrap {
	return projects.QueryerWrap{Q: w.q}
}

func (w queryerWrap) Columns() columns.QueryerWrap {
	return columns.QueryerWrap{Q: w.q}
}

func (w queryerWrap) Tasks() tasks.QueryerWrap {
	return tasks.QueryerWrap{Q: w.q}
}

func (w queryerWrap) Comments() comments.QueryerWrap {
	return comments.QueryerWrap{Q: w.q}
}

func QueryWithTX(tx TX) queryerWrap {
	return queryerWrap{q: tx}
}

func Query() queryerWrap {
	return queryerWrap{q: pool}
}

func Begin() (TX, error) {
	return pool.Begin(context.Background())
}

func Commit(tx TX) error {
	return tx.Commit(context.Background())
}

func Rollback(tx TX) {
	if err := tx.Rollback(context.Background()); err != nil {
		logger.Zap.Error("error while rollback db transaction", zap.Error(err))
	}
}

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
