package db

import (
	"encoding/base64"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

var defaultPoolConf = Config{
	Host:     "host",
	Port:     36,
	Name:     "db",
	User:     "user",
	Password: "password",

	ConnectionsPoolSize: 2,
}

func TestConfig_ToConnPoolConfig(t *testing.T) {
	assert := assert.New(t)

	c := defaultPoolConf
	actual, err := c.toConnPoolConfig()

	assert.Nil(err)
	assert.Equal("host", actual.ConnConfig.Host)
	assert.Equal(uint16(36), actual.ConnConfig.Port)
	assert.Equal("db", actual.ConnConfig.Database)
	assert.Equal("user", actual.ConnConfig.User)
	expectedPassword := base64.URLEncoding.EncodeToString([]byte(c.Password))
	assert.Equal(expectedPassword, actual.ConnConfig.Password)

	assert.Equal(int32(2), actual.MinConns)
}

func TestConfig_ToConnPoolConfig_WhenInvalidPort_ExpectError(t *testing.T) {
	assert := assert.New(t)

	c := defaultPoolConf
	c.Port = 0

	_, err := c.toConnPoolConfig()

	_, ok := err.(*pgconn.ParseConfigError)
	assert.True(ok)
}

func TestConfig_ToConnPoolConfig_UrlEncodesPassword(t *testing.T) {
	assert := assert.New(t)

	c := defaultPoolConf
	c.Password = "zefpoi*${oiz}"

	conf, err := c.toConnPoolConfig()

	assert.Nil(err)
	expectedConnString := "postgresql://user:emVmcG9pKiR7b2l6fQ==@host:36/db?pool_min_conns=2"
	assert.Equal(expectedConnString, conf.ConnString())
}
