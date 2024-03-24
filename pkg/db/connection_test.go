package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestConnection_ConnectUsesConnectionFunc(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultConnectionFunc)

	conf := Config{
		Host: "some-host",
	}
	conn := NewConnection(conf)

	called := false
	var actualConf pgx.ConnPoolConfig

	mockConnFunc := func(config pgx.ConnPoolConfig) (p *pgx.ConnPool, err error) {
		called = true
		actualConf = config
		return nil, err
	}
	pgxConnectionFunc = mockConnFunc

	conn.Connect()

	assert.True(called)
	expected := conf.toConnPoolConfig()
	assert.Equal(expected, actualConf)
}

var errDefault = fmt.Errorf("some error")

func TestConnection_ConnectPropagatesError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultConnectionFunc)

	conf := Config{
		Host: "some-host",
	}
	conn := NewConnection(conf)

	mockConnFunc := func(config pgx.ConnPoolConfig) (p *pgx.ConnPool, err error) {
		return nil, errDefault
	}
	pgxConnectionFunc = mockConnFunc

	err := conn.Connect()

	assert.Equal(errDefault, err)
}

type mockDbConnection struct {
	closeCalled int
	queryCalled int
	execCalled  int

	sqlQuery string
	args     []interface{}

	err error
}

func (m *mockDbConnection) Close() {
	m.closeCalled++
}

func (m *mockDbConnection) QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error) {
	m.queryCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return nil, m.err
}

func (m *mockDbConnection) ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error) {
	m.execCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return pgx.CommandTag(""), m.err
}

func TestConnection_CloseReleasesTheDbConnection(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	conn.Close()

	assert.Equal(1, m.closeCalled)
}

func resetDefaultConnectionFunc() {
	pgxConnectionFunc = pgx.NewConnPool
}

const exampleSqlQuery = "select * from table"

func TestConnection_Query_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	conn.Query(context.Background(), exampleSqlQuery)

	assert.Equal(1, m.queryCalled)
}

func TestConnection_Query_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	conn.Query(context.Background(), exampleSqlQuery)

	assert.Equal(exampleSqlQuery, m.sqlQuery)
}

func TestConnection_Query_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	conn.Query(context.Background(), exampleSqlQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.args)
}

func TestConnection_Query_DoesNotCreateError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	res := conn.Query(context.Background(), exampleSqlQuery)

	assert.Nil(res.Err())
}

func TestConnection_Query_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{
		err: errDefault,
	}
	conn := connectionImpl{
		pool: m,
	}

	res := conn.Query(context.Background(), exampleSqlQuery)

	assert.Equal(errDefault, res.Err())
}

const exampleExecQuery = "insert into table values('1')"

func TestConnection_Exec_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	conn.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, m.execCalled)
}

func TestConnection_Exec_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	conn.Exec(context.Background(), exampleExecQuery)

	assert.Equal(exampleExecQuery, m.sqlQuery)
}

func TestConnection_Exec_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	conn.Exec(context.Background(), exampleExecQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.args)
}

func TestConnection_Exec_DoesNotCreateError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{}
	conn := connectionImpl{
		pool: m,
	}

	_, err := conn.Exec(context.Background(), exampleExecQuery)

	assert.Nil(err)
}

func TestConnection_Exec_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockDbConnection{
		err: errDefault,
	}
	conn := connectionImpl{
		pool: m,
	}

	_, err := conn.Exec(context.Background(), exampleExecQuery)

	assert.Equal(errDefault, err)
}
