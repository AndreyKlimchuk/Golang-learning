package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AndreyKlimchuk/golang-learning/homework4/logger"
	"go.uber.org/zap"

	"github.com/AndreyKlimchuk/golang-learning/homework4/db/columns"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/comments"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/projects"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/tasks"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TX pgx.Tx

type queryerWrap common.QueryerWrap

var databaseURL string // = "postgres://gorello:12345@localhost:5432/gorello"

//const migrationsSourceUrl = "/db/migrations"
var pool *pgxpool.Pool

func Init() (err error) {
	databaseURL = os.Getenv("DATABASE_URL")
	if err := ApplyMigrationsUp(); err != nil {
		return fmt.Errorf("cannot apply up migrations: %w", err)
	}
	pool, err = pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("cannot connect pool: %w", err)
	}
	return nil
}

func ApplyMigrationsUp() error {
	return doMigrations(func(migrations *migrate.Migrate) error {
		return migrations.Up()
	})
}

func ApplyMigrationsDown() error {
	return doMigrations(func(migrations *migrate.Migrate) error {
		return migrations.Down()
	})
}

func doMigrations(do func(*migrate.Migrate) error) error {
	packageDir := getPackageDir()
	migrations, err := migrate.New("file://"+packageDir+"/migrations", databaseURL)
	if err != nil {
		return err
	}
	defer migrations.Close()
	if err := do(migrations); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func getPackageDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
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
