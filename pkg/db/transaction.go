package db

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/jackc/pgx"
)

type Transaction interface {
	Close(ctx context.Context) error

	Query(ctx context.Context, sql string, arguments ...interface{}) Rows
	Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error)
}

type pgxDbConnection interface {
	Close() error
	BeginEx(ctx context.Context, opts *pgx.TxOptions) (*pgx.Tx, error)
}

type pgxDbTransaction interface {
	RollbackEx(ctx context.Context) error
	CommitEx(ctx context.Context) error

	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error)
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (commandTag pgx.CommandTag, err error)

	Err() error
}

type transactionImpl struct {
	conn pgxDbConnection
	tx   pgxDbTransaction
	err  error
}

func defaultCreateTransactionFromPgxConnection(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
	return conn.BeginEx(ctx, nil)
}

var createTransactionFromPgxConn = defaultCreateTransactionFromPgxConnection

func newTransactionFromPool(ctx context.Context, pool pgxDbConnectionPool) (Transaction, error) {
	conn, err := pool.AcquireEx(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := createTransactionFromPgxConn(ctx, conn)
	if err != nil {
		return nil, err
	}

	transactionImpl := transactionImpl{
		conn: conn,
		tx:   tx,
	}

	return &transactionImpl, nil
}

func (t *transactionImpl) Close(ctx context.Context) error {
	var err error

	if t.err != nil {
		err = t.tx.RollbackEx(ctx)
	} else {
		err = t.tx.CommitEx(ctx)
	}

	err2 := t.tx.Err()
	var err3 error
	if t.conn != nil {
		// TODO: Should call t.pool.Release(t.conn)
		err3 = t.conn.Close()
	}

	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	return err3
}

func (t *transactionImpl) Query(ctx context.Context, sql string, arguments ...interface{}) Rows {
	log := middleware.GetLoggerFromContext(ctx)
	log.Debugf("Query: %s (%d)", sql, len(arguments))

	rows, err := t.tx.QueryEx(ctx, sql, nil, arguments...)
	t.updateErrorStatus(err)
	return newRows(rows, err)
}

func (t *transactionImpl) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	log := middleware.GetLoggerFromContext(ctx)
	log.Debugf("Exec: %s (%d)", sql, len(arguments))

	tag, err := t.tx.ExecEx(ctx, sql, nil, arguments...)
	t.updateErrorStatus(err)
	return int(tag.RowsAffected()), err
}

func (t *transactionImpl) updateErrorStatus(err error) {
	if err != nil {
		t.err = err
	}
}
