package db

import (
	"net/url"
	"testing"
	"time"

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

func TestUnit_Config_ToConnPoolConfig(t *testing.T) {
	assert := assert.New(t)

	c := defaultPoolConf
	actual, err := c.toConnPoolConfig()

	assert.Nil(err)
	assert.Equal("host", actual.ConnConfig.Host)
	assert.Equal(uint16(36), actual.ConnConfig.Port)
	assert.Equal("db", actual.ConnConfig.Database)
	assert.Equal("user", actual.ConnConfig.User)
	expectedPassword := url.QueryEscape(c.Password)
	assert.Equal(expectedPassword, actual.ConnConfig.Password)

	assert.Equal(int32(2), actual.MinConns)
	assert.Equal(time.Duration(0), actual.ConnConfig.ConnectTimeout)
}

func TestUnit_Config_ToConnPoolConfig_WithConnectTimeout(t *testing.T) {
	assert := assert.New(t)

	c := defaultPoolConf
	c.ConnectTimeout = 2 * time.Second
	actual, err := c.toConnPoolConfig()

	assert.Nil(err)

	assert.Equal(2*time.Second, actual.ConnConfig.ConnectTimeout)
}

func TestUnit_Config_ToConnPoolConfig_WhenConnectTimeoutNotAFullSecond_ExpectRounded(t *testing.T) {
	assert := assert.New(t)

	c := defaultPoolConf
	c.ConnectTimeout = 3*time.Second + 720*time.Millisecond
	actual, err := c.toConnPoolConfig()

	assert.Nil(err)

	assert.Equal(3*time.Second, actual.ConnConfig.ConnectTimeout)
}

func TestUnit_Config_ToConnPoolConfig_WhenInvalidPort_ExpectError(t *testing.T) {
	assert := assert.New(t)

	c := defaultPoolConf
	c.Port = 0

	_, err := c.toConnPoolConfig()

	_, ok := err.(*pgconn.ParseConfigError)
	assert.True(ok)
}

func TestUnit_Config_ToConnPoolConfig_UrlEncodesPassword(t *testing.T) {
	assert := assert.New(t)

	c := defaultPoolConf
	c.Password = "zefpoi*${oiz}"

	conf, err := c.toConnPoolConfig()

	assert.Nil(err)
	expectedConnString := "postgresql://user:zefpoi%2A%24%7Boiz%7D@host:36/db?pool_min_conns=2"
	assert.Equal(expectedConnString, conf.ConnString())
}
