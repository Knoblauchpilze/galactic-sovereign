package db

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/jackc/pgx"
)

type transactionImpl struct {
	tx  pgxTransaction
	err error
}

func (t *transactionImpl) Close(ctx context.Context) {
	var err error

	if t.err != nil {
		err = t.tx.RollbackEx(ctx)
	} else {
		err = t.tx.CommitEx(ctx)
	}

	if err != nil {
		middleware.GetLoggerFromContext(ctx).Warnf("Failed to finalize transaction: %v", err)
	}

	if err := t.tx.Err(); err != nil && err != pgx.ErrTxClosed {
		middleware.GetLoggerFromContext(ctx).Warnf("Transaction ended in error state: %v", err)
	}
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
