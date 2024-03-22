package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ToConnPoolConfig(t *testing.T) {
	assert := assert.New(t)

	c := Config{
		Host:     "host",
		Port:     36,
		Name:     "db",
		User:     "user",
		Password: "password",

		ConnectionsPoolSize: 2,
	}

	actual := c.toConnPoolConfig()

	assert.Equal("host", actual.Host)
	assert.Equal(uint16(36), actual.Port)
	assert.Equal("db", actual.Database)
	assert.Equal("user", actual.User)
	assert.Equal("password", actual.Password)

	assert.Equal(2, actual.MaxConnections)
	assert.Equal(time.Duration(0), actual.AcquireTimeout)
}
