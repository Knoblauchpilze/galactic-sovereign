package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Connection interface {
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	BeginTx(ctx context.Context) (Transaction, error)
	Exec(ctx context.Context, sql string, arguments ...any) (int64, error)
	Close(ctx context.Context)

	Ping(ctx context.Context) error
}
