package internal

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig_LeavesServerUnchanged(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := rest.Config{
		BasePath:  "/v1",
		Port:      uint16(80),
		RateLimit: 10,
	}
	assert.Equal(expected, conf.Server)
}

func TestDefaultConfig_LeavesApiKeyUnchanged(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := service.ApiConfig{
		ApiKeyValidity: 3 * time.Hour,
	}
	assert.Equal(expected, conf.ApiKey)
}

func TestDefaultConfig_ReplacesDatabaseConfiguration(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := db.Config{
		Host:                "172.17.0.1",
		Port:                5432,
		Name:                "db_stellar_dominion",
		ConnectionsPoolSize: 1,
		LogLevel:            log.DEBUG,
	}
	assert.Equal(expected, conf.Database)
}
