package internal

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/rest"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestUnit_DefaultConfig_SetsCorrectPrefix(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := rest.Config{
		BasePath:  "/v1",
		Prefix:    "/galactic-sovereign",
		Port:      uint16(80),
		RateLimit: 10,
	}
	assert.Equal(expected, conf.Server)
}

func TestUnit_DefaultConfig_LeavesApiKeyUnchanged(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := service.ApiConfig{
		ApiKeyValidity: 3 * time.Hour,
	}
	assert.Equal(expected, conf.ApiKey)
}

func TestUnit_DefaultConfig_ReplacesDatabaseConfiguration(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := db.Config{
		Host:                "172.17.0.1",
		Port:                5432,
		Name:                "db_galactic_sovereign",
		ConnectionsPoolSize: 1,
		ConnectTimeout:      2 * time.Second,
		LogLevel:            log.DEBUG,
	}
	assert.Equal(expected, conf.Database)
}
