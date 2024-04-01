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

type mockPgxDbTransaction struct {
	queryCalled    int
	execCalled     int
	rollbackCalled int
	commitCalled   int

	sqlQuery  string
	arguments []interface{}

	tag         pgx.CommandTag
	rollbackErr error
	commitErr   error
	err         error
}

func TestTransaction_Query_DelegatesToTransaction(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	tx.Query(context.Background(), exampleSqlQuery)

	assert.Equal(1, mt.queryCalled)
}

func TestTransaction_Query_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	tx.Query(context.Background(), exampleSqlQuery)

	assert.Equal(exampleSqlQuery, mt.sqlQuery)
}

func TestTransaction_Query_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	tx.Query(context.Background(), exampleSqlQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, mt.arguments)
}

func TestTransaction_Query_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{
		err: errDefault,
	}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	actual := tx.Query(context.Background(), exampleSqlQuery)

	assert.Equal(errDefault, actual.Err())
}

func TestTransaction_Exec_DelegatesToTransaction(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	tx.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, mt.execCalled)
}

func TestTransaction_Exec_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	tx.Exec(context.Background(), exampleExecQuery)

	assert.Equal(exampleExecQuery, mt.sqlQuery)
}

func TestTransaction_Exec_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	tx.Exec(context.Background(), exampleExecQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, mt.arguments)
}

func TestTransaction_Exec_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{
		err: errDefault,
	}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	_, err := tx.Exec(context.Background(), exampleExecQuery)

	assert.Equal(errDefault, err)
}

func TestTransaction_Exec_PropagatesCommandTag(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionCreationFunc)

	mp := mockPgxDbConnectionPool{}
	mt := mockPgxDbTransaction{
		tag: pgx.CommandTag("INSERT 0 1"),
	}
	createTransactionFromPgxConn = func(ctx context.Context, conn *pgx.Conn) (pgxDbTransaction, error) {
		return &mt, nil
	}

	tx, _ := newTransactionFromPool(context.Background(), &mp)

	actual, _ := tx.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, actual)
}

func TestTransaction_Close_WhenError_CallsRollback(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxDbTransaction{}
	tx := transactionImpl{
		tx: &mt,

		err: errDefault,
	}

	tx.Close(context.Background())

	assert.Equal(1, mt.rollbackCalled)
}

func TestTransaction_Close_ReturnsRollbackError(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxDbTransaction{
		rollbackErr: errDefault,
	}
	tx := transactionImpl{
		tx: &mt,

		err: errDefault,
	}

	actual := tx.Close(context.Background())

	assert.Equal(errDefault, actual)
}

func TestTransaction_Close_WhenNoError_CallsCommit(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxDbTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	tx.Close(context.Background())

	assert.Equal(1, mt.commitCalled)
}

func TestTransaction_Close_ReturnsCommitError(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxDbTransaction{
		commitErr: errDefault,
	}
	tx := transactionImpl{
		tx: &mt,
	}

	actual := tx.Close(context.Background())

	assert.Equal(errDefault, actual)
}

func TestTransaction_Close_WhenCommitSucceeds_ReturnsGeneralErrorState(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxDbTransaction{
		err: errDefault,
	}
	tx := transactionImpl{
		tx: &mt,
	}

	actual := tx.Close(context.Background())

	assert.Equal(errDefault, actual)
}

func TestTransaction_Close_NominalCase(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxDbTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	actual := tx.Close(context.Background())

	assert.Nil(actual)
}

func resetConnectionCreationFunc() {
	createTransactionFromPgxConn = defaultCreateTransactionFromPgxConnection
}

func (m *mockPgxDbTransaction) RollbackEx(ctx context.Context) error {
	m.rollbackCalled++
	return m.rollbackErr
}

func (m *mockPgxDbTransaction) CommitEx(ctx context.Context) error {
	m.commitCalled++
	return m.commitErr
}

func (m *mockPgxDbTransaction) QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error) {
	m.queryCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return nil, m.err
}

func (m *mockPgxDbTransaction) ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (commandTag pgx.CommandTag, err error) {
	m.execCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return m.tag, m.err
}

func (m *mockPgxDbTransaction) Err() error {
	return m.err
}
