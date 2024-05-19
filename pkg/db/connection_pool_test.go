package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

func TestConnectionPool_NewConnectionPool_UsesProvidedConfig(t *testing.T) {
	assert := assert.New(t)

	p := NewConnectionPool(defaultPoolConf)

	actual, ok := p.(*connectionPoolImpl)

	assert.True(ok)
	assert.Equal(defaultPoolConf, actual.config)
}

func TestConnectionPool_ConnectUsesConnectionFunc(t *testing.T) {
	assert := assert.New(t)

	called := false

	mockConnFunc := func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		called = true
		return nil, nil
	}

	p := newConnectionPool(defaultPoolConf, mockConnFunc)
	p.Connect(context.Background())

	assert.True(called)
}

func TestConnectionPool_ConnectToExpectedDatabase(t *testing.T) {
	assert := assert.New(t)

	var actualConf *pgxpool.Config

	mockConnFunc := func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		actualConf = config
		return nil, nil
	}

	p := newConnectionPool(defaultPoolConf, mockConnFunc)
	err := p.Connect(context.Background())

	expected, _ := defaultPoolConf.toConnPoolConfig()

	assert.Nil(err)
	assert.Equal(expected.ConnString(), actualConf.ConnString())
}

func TestConnectionPool_ConnectPropagatesConversionError(t *testing.T) {
	assert := assert.New(t)

	mockConnFunc := func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		return nil, nil
	}

	conf := defaultPoolConf
	conf.Port = 0

	p := newConnectionPool(conf, mockConnFunc)
	err := p.Connect(context.Background())

	actual, ok := err.(*pgconn.ParseConfigError)
	assert.True(ok)
	assert.Contains(actual.Error(), "invalid port (outside range)")
}

func TestConnectionPool_ConnectPropagatesConnectionError(t *testing.T) {
	assert := assert.New(t)

	mockConnFunc := func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		return nil, errDefault
	}

	p := newConnectionPool(defaultPoolConf, mockConnFunc)
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

func TestConnectionPool_Close_ReleasesTheDbConnection(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Close()

	assert.Equal(1, m.closeCalled)
}

func TestConnectionPool_Ping_DelegatesToPool(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{}
	p := connectionPoolImpl{
		pool: m,
	}

	p.Ping(context.Background())

	assert.Equal(1, m.pingCalled)
}

func TestConnectionPool_Ping_PropagatesPoolError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPgxConnectionPool{
		err: errDefault,
	}
	p := connectionPoolImpl{
		pool: m,
	}

	err := p.Ping(context.Background())

	assert.Equal(errDefault, err)
}

func TestConnectionPool_Ping_WhenPoolReturnsConnectionError_IndicatesFailure(t *testing.T) {
	assert := assert.New(t)

	connErr := pgconn.ConnectError{
		Config: &pgconn.Config{},
	}

	m := &mockPgxConnectionPool{
		err: &connErr,
	}
	p := connectionPoolImpl{
		pool: m,
	}

	err := p.Ping(context.Background())

	assert.True(errors.IsErrorWithCode(err, DatabasePingFailed))
}

func TestConnectionPool_Ping_WhenPoolReturnsConnectionError_WrapsIntoError(t *testing.T) {
	assert := assert.New(t)

	connErr := pgconn.ConnectError{
		Config: &pgconn.Config{},
	}

	m := &mockPgxConnectionPool{
		err: &connErr,
	}
	p := connectionPoolImpl{
		pool: m,
	}

	err := p.Ping(context.Background())

	actual := errors.Unwrap(err)
	assert.NotNil(actual)
	actualCause, ok := actual.(*pgconn.ConnectError)
	assert.True(ok)
	assert.Equal(&connErr, actualCause)
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
