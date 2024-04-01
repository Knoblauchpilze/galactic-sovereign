package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_AcquiresNewConnectionFromPool(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return nil, nil
	}
	m := mockPgxDbConnectionPool{}

	newTransactionFromPool(context.Background(), &m)

	assert.Equal(1, m.acquireCalled)
}

func TestTransaction_WhenAcquiringTransactionFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	m := mockPgxDbConnectionPool{
		err: errDefault,
	}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return nil, nil
	}

	_, err := newTransactionFromPool(context.Background(), &m)

	assert.Equal(errDefault, err)
}

func TestTransaction_StartsTransactionOnConnection(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	called := false

	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		called = true
		return nil, nil
	}
	m := mockPgxDbConnectionPool{}

	newTransactionFromPool(context.Background(), &m)

	assert.True(called)
}

func TestTransaction_UsesTransactionImplementation(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return nil, nil
	}
	m := mockPgxDbConnectionPool{}

	actual, err := newTransactionFromPool(context.Background(), &m)

	assert.Nil(err)
	assert.IsType(&transactionImpl{}, actual)
}

func TestTransaction_WhenStartingTransactionFails_ExpectError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return nil, errDefault
	}
	m := mockPgxDbConnectionPool{}

	_, err := newTransactionFromPool(context.Background(), &m)

	assert.Equal(errDefault, err)
}

type mockPgxDbConnection struct {
	pgxDbConnection
}

type mockPgxDbTransaction struct {
	pgxDbTransaction
}

func resetConnectionCreationFunc() {
	createTransactionFromPgxConn = defaultCreateTransactionFromPgxConnection
}
