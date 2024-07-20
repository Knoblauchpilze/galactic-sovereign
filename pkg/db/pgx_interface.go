package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// An implicit interface representing what pgxpool.Pool is offering.
// Only the needed methods are exposed.
type pgxConnectionPool interface {
	Close()

	Ping(ctx context.Context) error

	Begin(ctx context.Context) (pgx.Tx, error)

	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// The public interface of what we expose in a connection pool.
// This is agnostic to pgx.
type ConnectionPool interface {
	Connect(ctx context.Context) error
	Close()

	Ping(ctx context.Context) error

	StartTransaction(ctx context.Context) (Transaction, error)

	Query(ctx context.Context, sql string, arguments ...any) Rows
	Exec(ctx context.Context, sql string, arguments ...any) (int, error)
}

// The public interface of what we expose in a transaction.
// This is agnostic to pgx.
type Transaction interface {
	Close(ctx context.Context)

	Query(ctx context.Context, sql string, arguments ...any) Rows
	Exec(ctx context.Context, sql string, arguments ...any) (int, error)
}

// The public interface of what we expose in sql rows.
// This is agnostic to pgx.
type Rows interface {
	Err() error

	GetSingleValue(parser RowParser) error
	GetAll(parser RowParser) error
}
