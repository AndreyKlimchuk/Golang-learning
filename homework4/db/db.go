package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"

	"github.com/AndreyKlimchuk/golang-learning/homework4/logger"
	"go.uber.org/zap"

	"github.com/AndreyKlimchuk/golang-learning/homework4/db/columns"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/comments"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/projects"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/tasks"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TX pgx.Tx

var pool *pgxpool.Pool

type queryerWrap common.QueryerWrap

func Init(conn string) (err error) {
	pool, err = pgxpool.Connect(context.Background(), conn)
	if err != nil {
		return err
	}
	return nil
}

func (w queryerWrap) Projects() projects.QueryerWrap {
	return projects.QueryerWrap(w)
}

func (w queryerWrap) Columns() columns.QueryerWrap {
	return columns.QueryerWrap(w)
}

func (w queryerWrap) Tasks() tasks.QueryerWrap {
	return tasks.QueryerWrap(w)
}

func (w queryerWrap) Comments() comments.QueryerWrap {
	return comments.QueryerWrap(w)
}

func QueryWithTX(tx TX) queryerWrap {
	return queryerWrap{Q: tx}
}

func Query() queryerWrap {
	return queryerWrap{Q: pool}
}

func Begin() (TX, error) {
	return pool.Begin(context.Background())
}

func Commit(tx TX) error {
	return tx.Commit(context.Background())
}

func Rollback(tx TX) {
	if err := tx.Rollback(context.Background()); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		logger.Zap.Error("error while rollback db transaction", zap.Error(err))
	}
}
