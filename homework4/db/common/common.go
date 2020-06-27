package common

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type QueryerWrap struct {
	Q Queryer
}

var ErrNoAffectedRows = errors.New("no affected rows")
var ErrNoRows = pgx.ErrNoRows

type Queryer interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func IsNoRowsError(err error) bool {
	// TODO: check correctness
	return errors.Is(err, ErrNoRows) || errors.Is(err, ErrNoAffectedRows)
}

func ErrorIfNoAffectedRows(ct pgconn.CommandTag, err error) error {
	if err != nil {
		return err
	} else if ct.RowsAffected() == 0 {
		return ErrNoAffectedRows
	}
	return nil
}
