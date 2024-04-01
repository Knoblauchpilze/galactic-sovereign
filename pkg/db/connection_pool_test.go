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

type mockDbConnectionPool struct {
	closeCalled int
	queryCalled int
	execCalled  int

	sqlQuery  string
	arguments []interface{}

	err error
}

func TestConnectionPool_CloseReleasesTheDbConnection(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Close()

	assert.Equal(1, m.closeCalled)
}

const exampleSqlQuery = "select * from table"

func TestConnectionPool_Query_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Query(context.Background(), exampleSqlQuery)

	assert.Equal(1, m.queryCalled)
}

func TestConnectionPool_Query_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Query(context.Background(), exampleSqlQuery)

	assert.Equal(exampleSqlQuery, m.sqlQuery)
}

func TestConnectionPool_Query_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Query(context.Background(), exampleSqlQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.arguments)
}

func TestConnectionPool_Query_DoesNotCreateError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	res := connPool.Query(context.Background(), exampleSqlQuery)

	assert.Nil(res.Err())
}

func TestConnectionPool_Query_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{
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

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, m.execCalled)
}

func TestConnectionPool_Exec_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Exec(context.Background(), exampleExecQuery)

	assert.Equal(exampleExecQuery, m.sqlQuery)
}

func TestConnection_Exec_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	connPool.Exec(context.Background(), exampleExecQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.arguments)
}

func TestConnectionPool_Exec_DoesNotCreateError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{}
	connPool := connectionPoolImpl{
		pool: m,
	}

	_, err := connPool.Exec(context.Background(), exampleExecQuery)

	assert.Nil(err)
}

func TestConnectionPool_Exec_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnectionPool{
		err: errDefault,
	}
	connPool := connectionPoolImpl{
		pool: m,
	}

	_, err := connPool.Exec(context.Background(), exampleExecQuery)

	assert.Equal(errDefault, err)
}

func (m *mockDbConnectionPool) Close() {
	m.closeCalled++
}

func (m *mockDbConnectionPool) AcquireEx(ctx context.Context) (*pgx.Conn, error) {
	return nil, nil
}

func (m *mockDbConnectionPool) QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error) {
	m.queryCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return nil, m.err
}

func (m *mockDbConnectionPool) ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error) {
	m.execCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return pgx.CommandTag(""), m.err
}
