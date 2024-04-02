package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

func TestConnectionPool_ConnectUsesConnectionFunc(t *testing.T) {
	assert := assert.New(t)

	called := false
	var actualConf pgx.ConnPoolConfig

	mockConnFunc := func(config pgx.ConnPoolConfig) (*pgx.ConnPool, error) {
		called = true
		actualConf = config
		return nil, nil
	}

	conf := Config{
		Host: "some-host",
	}

	p := newConnectionPool(conf, mockConnFunc)
	p.Connect()

	assert.True(called)
	expected := conf.toConnPoolConfig()
	assert.Equal(expected, actualConf)
}

func TestConnectionPool_ConnectPropagatesError(t *testing.T) {
	assert := assert.New(t)

	mockConnFunc := func(config pgx.ConnPoolConfig) (*pgx.ConnPool, error) {
		return nil, errDefault
	}

	p := newConnectionPool(Config{}, mockConnFunc)
	err := p.Connect()

	assert.Equal(errDefault, err)
}

type mockPgxConnectionPool struct {
	closeCalled int
	beginCalled int
	queryCalled int
	execCalled  int

	sqlQuery  string
	arguments []interface{}

	tx  *pgx.Tx
	err error
}

func TestConnectionPool_CloseReleasesTheDbConnection(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Close()

	assert.Equal(1, m.closeCalled)
}

func TestConnectionPool_StartTransaction_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.StartTransaction(context.Background())

	assert.Equal(1, m.beginCalled)
}

func TestConnectionPool_StartTransaction_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{
		err: errDefault,
	}
	p := connectionPoolImpl{
		pool: m,
	}

	_, err := p.StartTransaction(context.Background())

	assert.Equal(errDefault, err)
}

func TestConnectionPool_StartTransaction_CreatesTransaction(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{
		tx: &pgx.Tx{},
	}
	p := connectionPoolImpl{
		pool: m,
	}

	tx, err := p.StartTransaction(context.Background())

	assert.Nil(err)
	actual, ok := tx.(*transactionImpl)
	assert.True(ok)
	assert.Equal(m.tx, actual.tx)
}

const exampleSqlQuery = "select * from table"

func TestConnectionPool_Query_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Query(context.Background(), exampleSqlQuery)

	assert.Equal(1, m.queryCalled)
}

func TestConnectionPool_Query_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Query(context.Background(), exampleSqlQuery)

	assert.Equal(exampleSqlQuery, m.sqlQuery)
}

func TestConnectionPool_Query_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Query(context.Background(), exampleSqlQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.arguments)
}

func TestConnectionPool_Query_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{
		err: errDefault,
	}
	p := connectionPoolImpl{
		pool: m,
	}

	res := p.Query(context.Background(), exampleSqlQuery)

	assert.Equal(errDefault, res.Err())
}

const exampleExecQuery = "insert into table values('1')"

func TestConnectionPool_Exec_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, m.execCalled)
}

func TestConnectionPool_Exec_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Exec(context.Background(), exampleExecQuery)

	assert.Equal(exampleExecQuery, m.sqlQuery)
}

func TestConnection_Exec_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Exec(context.Background(), exampleExecQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.arguments)
}

func TestConnectionPool_Exec_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{
		err: errDefault,
	}
	p := connectionPoolImpl{
		pool: m,
	}

	_, err := p.Exec(context.Background(), exampleExecQuery)

	assert.Equal(errDefault, err)
}

func (m *mockPgxConnectionPool) Close() {
	m.closeCalled++
}

func (m *mockPgxConnectionPool) BeginEx(ctx context.Context, txOptions *pgx.TxOptions) (*pgx.Tx, error) {
	m.beginCalled++
	return m.tx, m.err
}

func (m *mockPgxConnectionPool) QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error) {
	m.queryCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return nil, m.err
}

func (m *mockPgxConnectionPool) ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error) {
	m.execCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return pgx.CommandTag(""), m.err
}
