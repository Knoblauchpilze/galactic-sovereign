package internal

import (
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig_Server(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	expected := rest.Config{
		Endpoint: "/v1/users/",
		Port:     uint16(60000),
	}
	assert.Equal(expected, conf.Server)
}

func TestDefaultConfig_Database_AssumesLocalhost(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	assert.Equal("localhost", conf.Database.Host)
}

func TestDefaultConfig_Database_AssumesPort5432(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	assert.Equal(uint16(5432), conf.Database.Port)
}

func TestDefaultConfig_Database_AssumesPoolSizeOfOne(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	assert.Equal(uint(1), conf.Database.ConnectionsPoolSize)
}

func TestDefaultConfig_Database_DoesNotDefineDatabase(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	assert.Equal("", conf.Database.Name)
}

func TestDefaultConfig_Database_DoesNotDefineUser(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	assert.Equal("", conf.Database.User)
}

func TestDefaultConfig_Database_DoesNotDefinePassword(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	assert.Equal("", conf.Database.Password)
}
