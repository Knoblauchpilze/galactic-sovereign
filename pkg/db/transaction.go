package db

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type transactionImpl struct {
	tx  pgx.Tx
	err error
}

func (t *transactionImpl) Close(ctx context.Context) {
	var err error

	if t.err != nil {
		err = t.tx.Rollback(ctx)
	} else {
		err = t.tx.Commit(ctx)
	}

	if err != nil && err != pgx.ErrTxClosed {
		logger.GetRequestLogger(ctx).Warnf("Transaction ended in error state: %v", err)
	}
}

func (t *transactionImpl) Query(ctx context.Context, sql string, arguments ...interface{}) Rows {
	rows, err := t.tx.Query(ctx, sql, arguments...)
	t.updateErrorStatus(err)
	return newRows(rows, err)
}

func (t *transactionImpl) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	tag, err := t.tx.Exec(ctx, sql, arguments...)
	t.updateErrorStatus(err)
	return int(tag.RowsAffected()), err
}

func (t *transactionImpl) updateErrorStatus(err error) {
	if err != nil {
		t.err = err
	}
}
