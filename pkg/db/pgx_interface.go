package db

import (
	"context"

	"github.com/jackc/pgx"
)

// An implicit interface representing what pgx.ConnPool is offering.
// Only the needed methods are exposed.
type pgxConnectionPool interface {
	Close()

	BeginEx(ctx context.Context, txOptions *pgx.TxOptions) (*pgx.Tx, error)

	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error)
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error)
}

// An implicit interface representing what pgx.Tx is offering.
// Only the needed methods are exposed.
type pgxTransaction interface {
	RollbackEx(ctx context.Context) error
	CommitEx(ctx context.Context) error

	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error)
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (commandTag pgx.CommandTag, err error)

	Err() error
}

// An implicit interface representing what pgx.Rows is offering.
// Only the needed methods are exposed.
type pgxRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
}

// The public interface of what we expose in a connection pool.
// This is agnostic to pgx.
type ConnectionPool interface {
	Connect() error
	Close()

	StartTransaction(ctx context.Context) (Transaction, error)

	Query(ctx context.Context, sql string, arguments ...interface{}) Rows
	Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error)
}

// The public interface of what we expose in a transaction.
// This is agnostic to pgx.
type Transaction interface {
	Close(ctx context.Context)

	Query(ctx context.Context, sql string, arguments ...interface{}) Rows
	Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error)
}

// The public interface of what we expose in sql rows.
// This is agnostic to pgx.
type Rows interface {
	Err() error

	Close()

	GetSingleValue(parser RowParser) error
	GetAll(parser RowParser) error
}
