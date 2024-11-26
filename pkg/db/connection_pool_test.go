package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

func TestUnit_ConnectionPool_NewConnectionPool_UsesProvidedConfig(t *testing.T) {
	assert := assert.New(t)

	p := NewConnectionPool(defaultPoolConf, &mockLogger{})

	actual, ok := p.(*connectionPoolImpl)

	assert.True(ok)
	assert.Equal(defaultPoolConf, actual.config)
}

func TestUnit_ConnectionPool_ConnectUsesConnectionFunc(t *testing.T) {
	assert := assert.New(t)

	called := false

	mockConnFunc := func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		called = true
		return nil, nil
	}

	p := newConnectionPool(defaultPoolConf, &mockLogger{}, mockConnFunc)
	p.Connect(context.Background())

	assert.True(called)
}

func TestUnit_ConnectionPool_ConnectToExpectedDatabase(t *testing.T) {
	assert := assert.New(t)

	var actualConf *pgxpool.Config

	mockConnFunc := func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		actualConf = config
		return nil, nil
	}

	p := newConnectionPool(defaultPoolConf, &mockLogger{}, mockConnFunc)
	err := p.Connect(context.Background())

	expected, _ := defaultPoolConf.toConnPoolConfig(&mockLogger{})

	assert.Nil(err)
	assert.Equal(expected.ConnString(), actualConf.ConnString())
}

func TestUnit_ConnectionPool_ConnectPropagatesConversionError(t *testing.T) {
	assert := assert.New(t)

	mockConnFunc := func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		return nil, nil
	}

	conf := defaultPoolConf
	conf.Port = 0

	p := newConnectionPool(conf, &mockLogger{}, mockConnFunc)
	err := p.Connect(context.Background())

	actual, ok := err.(*pgconn.ParseConfigError)
	assert.True(ok)
	assert.Contains(actual.Error(), "invalid port (outside range)")
}

func TestUnit_ConnectionPool_ConnectPropagatesConnectionError(t *testing.T) {
	assert := assert.New(t)

	mockConnFunc := func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		return nil, errDefault
	}

	p := newConnectionPool(defaultPoolConf, &mockLogger{}, mockConnFunc)
	err := p.Connect(context.Background())

	assert.Equal(errDefault, err)
}

type mockPgxConnectionPool struct {
	closeCalled int
	pingCalled  int
	beginCalled int
	queryCalled int
	execCalled  int

	sqlQuery  string
	arguments []interface{}

	tx  pgx.Tx
	err error
}

func TestUnit_ConnectionPool_Close_ReleasesTheDbConnection(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		log:  &mockLogger{},
		pool: m,
	}

	p.Close()

	assert.Equal(1, m.closeCalled)
}

func TestUnit_ConnectionPool_Ping_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Ping(context.Background())

	assert.Equal(1, m.pingCalled)
}

func TestUnit_ConnectionPool_Ping_WhenPoolReturnsError_IndicatesFailure(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{
		err: errDefault,
	}
	p := connectionPoolImpl{
		pool: m,
	}

	err := p.Ping(context.Background())

	assert.True(errors.IsErrorWithCode(err, DatabasePingFailed))
}

func TestUnit_ConnectionPool_Ping_WhenPoolReturnsError_WrapsIntoError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{
		err: errDefault,
	}
	p := connectionPoolImpl{
		pool: m,
	}

	err := p.Ping(context.Background())

	actual := errors.Unwrap(err)
	assert.NotNil(actual)
	assert.Equal(errDefault, actual)
}

func TestUnit_ConnectionPool_StartTransaction_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.StartTransaction(context.Background())

	assert.Equal(1, m.beginCalled)
}

func TestUnit_ConnectionPool_StartTransaction_PropagatesPoolError(t *testing.T) {
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

func TestUnit_ConnectionPool_StartTransaction_CreatesTransaction(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{
		tx: &mockPgxTransaction{},
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

func TestUnit_ConnectionPool_Query_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Query(context.Background(), exampleSqlQuery)

	assert.Equal(1, m.queryCalled)
}

func TestUnit_ConnectionPool_Query_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Query(context.Background(), exampleSqlQuery)

	assert.Equal(exampleSqlQuery, m.sqlQuery)
}

func TestUnit_ConnectionPool_Query_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Query(context.Background(), exampleSqlQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.arguments)
}

func TestUnit_ConnectionPool_Query_PropagatesPoolError(t *testing.T) {
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

func TestUnit_ConnectionPool_Exec_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, m.execCalled)
}

func TestUnit_ConnectionPool_Exec_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Exec(context.Background(), exampleExecQuery)

	assert.Equal(exampleExecQuery, m.sqlQuery)
}

func TestUnit_Connection_Exec_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Exec(context.Background(), exampleExecQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, m.arguments)
}

func TestUnit_ConnectionPool_Exec_PropagatesPoolError(t *testing.T) {
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

func (m *mockPgxConnectionPool) Ping(ctx context.Context) error {
	m.pingCalled++
	return m.err
}

func (m *mockPgxConnectionPool) Begin(ctx context.Context) (pgx.Tx, error) {
	m.beginCalled++
	return m.tx, m.err
}

func (m *mockPgxConnectionPool) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	m.queryCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return nil, m.err
}

func (m *mockPgxConnectionPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.execCalled++
	m.sqlQuery = sql
	m.arguments = append(m.arguments, arguments...)
	return pgconn.NewCommandTag(""), m.err
}
