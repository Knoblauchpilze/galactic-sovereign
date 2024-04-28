package db

import (
	"testing"

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

	assert.Equal("host", actual.ConnConfig.Host)
	assert.Equal(uint16(36), actual.ConnConfig.Port)
	assert.Equal("db", actual.ConnConfig.Database)
	assert.Equal("user", actual.ConnConfig.User)
	assert.Equal("password", actual.ConnConfig.Password)

	assert.Equal(int32(2), actual.MinConns)
}
