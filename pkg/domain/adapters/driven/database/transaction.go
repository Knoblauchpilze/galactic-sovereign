package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Transaction interface {
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (int64, error)
	Rollback() error
	Close(ctx context.Context)
}
