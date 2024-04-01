package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestConnectionPool_ConnectUsesConnectionFunc(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionPoolFunc)

	conf := Config{
		Host: "some-host",
	}
	connPool := NewConnectionPool(conf)

	called := false
	var actualConf pgx.ConnPoolConfig

	mockConnFunc := func(config pgx.ConnPoolConfig) (p *pgx.ConnPool, err error) {
		called = true
		actualConf = config
		return nil, err
	}
	pgxConnectionFunc = mockConnFunc

	connPool.Connect()

	assert.True(called)
	expected := conf.toConnPoolConfig()
	assert.Equal(expected, actualConf)
}

var errDefault = fmt.Errorf("some error")

func TestConnectionPool_ConnectPropagatesError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConnectionPoolFunc)

	conf := Config{
		Host: "some-host",
	}
	connPool := NewConnectionPool(conf)

	mockConnFunc := func(config pgx.ConnPoolConfig) (p *pgx.ConnPool, err error) {
		return nil, errDefault
	}
	pgxConnectionFunc = mockConnFunc

	err := connPool.Connect()

	assert.Equal(errDefault, err)
}

func resetConnectionPoolFunc() {
	pgxConnectionFunc = pgx.NewConnPool
}

type mockPgxDbConnectionPool struct {
	closeCalled   int
	acquireCalled int
	queryCalled   int
	execCalled    int

	sqlQuery  string
	arguments []interface{}

	conn *pgx.Conn
	err  error
}

func TestConnectionPool_CloseReleasesTheDbConnection(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Close()

	assert.Equal(1, m.closeCalled)
}

func TestConnectionPool_BeginTransaction_CallsTransactionCreationFunc(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetTransactionFunc)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	called := false
	var actualPool pgxDbConnectionPool

	mockTransactionFunc := func(ctx context.Context, pool pgxDbConnectionPool) (Transaction, error) {
		called = true
		actualPool = pool
		return nil, nil
	}
	pgxTransactionFunc = mockTransactionFunc

	_, err := connPool.BeginTransaction(context.Background())

	assert.Nil(err)

	assert.True(called)
	assert.Equal(m, actualPool)
}

func resetTransactionFunc() {
	pgxTransactionFunc = newTransactionFromPool
}

const exampleSqlQuery = "select * from table"

func TestConnectionPool_Query_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Query(context.Background(), exampleSqlQuery)

	assert.Equal(1, m.queryCalled)
}

func TestConnectionPool_Query_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Query(context.Background(), exampleSqlQuery)

	assert.Equal(exampleSqlQuery, m.sqlQuery)
}

func TestConnectionPool_Query_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Query(context.Background(), exampleSqlQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.arguments)
}

func TestConnectionPool_Query_DoesNotCreateError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	res := connPool.Query(context.Background(), exampleSqlQuery)

	assert.Nil(res.Err())
}

func TestConnectionPool_Query_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{
		err: errDefault,
	}
	connPool := connectionPoolImpl{
		pool: m,
	}

	res := connPool.Query(context.Background(), exampleSqlQuery)

	assert.Equal(errDefault, res.Err())
}

const exampleExecQuery = "insert into table values('1')"

func TestConnectionPool_Exec_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, m.execCalled)
}

func TestConnectionPool_Exec_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Exec(context.Background(), exampleExecQuery)

	assert.Equal(exampleExecQuery, m.sqlQuery)
}

func TestConnection_Exec_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Exec(context.Background(), exampleExecQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.arguments)
}

func TestConnectionPool_Exec_DoesNotCreateError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	_, err := connPool.Exec(context.Background(), exampleExecQuery)

	assert.Nil(err)
}

func TestConnectionPool_Exec_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxDbConnectionPool{
		err: errDefault,
	}
	connPool := connectionPoolImpl{
		pool: m,
	}

	_, err := connPool.Exec(context.Background(), exampleExecQuery)

	assert.Equal(errDefault, err)
}

func (m *mockPgxDbConnectionPool) Close() {
	m.closeCalled++
}

func (m *mockPgxDbConnectionPool) AcquireEx(ctx context.Context) (*pgx.Conn, error) {
	m.acquireCalled++
	return m.conn, m.err
}

func (m *mockPgxDbConnectionPool) QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error) {
	m.queryCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return nil, m.err
}

func (m *mockPgxDbConnectionPool) ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error) {
	m.execCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return pgx.CommandTag(""), m.err
}
