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

func resetDefaultConnectionFunc() {
	pgxConnectionFunc = pgx.NewConnPool
}
