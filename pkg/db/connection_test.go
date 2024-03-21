package db

import (
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
	conn := New(conf)

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
	conn := New(conf)

	mockConnFunc := func(config pgx.ConnPoolConfig) (p *pgx.ConnPool, err error) {
		return nil, errDefault
	}
	pgxConnectionFunc = mockConnFunc

	err := conn.Connect()

	assert.Equal(errDefault, err)
}

type mockDbConnection struct {
	closeCalled int
}

func (m *mockDbConnection) Close() {
	m.closeCalled++
}

func (m *mockDbConnection) Query(sql string, args ...interface{}) (*pgx.Rows, error) {
	return nil, nil
}

func (m *mockDbConnection) Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error) {
	return pgx.CommandTag(""), nil
}

func TestConnection_CloseReleasesTheDbConnection(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultConnectionFunc)

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
