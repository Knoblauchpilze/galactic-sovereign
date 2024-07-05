package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/config"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

func TestDefaultConfig_ServerIsUnchanged(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	expected := rest.Config{
		BasePath:  "/v1",
		Port:      uint16(80),
		RateLimit: 10,
	}
	assert.Equal(expected, conf.Server)
}

func TestDefaultConfig_ServerDefinesTheDbName(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	expected := db.Config{
		Host:                "172.17.0.1",
		Port:                5432,
		Name:                "db_user_service",
		User:                "",
		Password:            "",
		ConnectionsPoolSize: 1,
		LogLevel:            log.DEBUG,
	}
	assert.Equal(expected, conf.Database)
}

func TestDefaultConfig_ApiKey(t *testing.T) {
	assert := assert.New(t)

	conf := defaultConf()

	expected := service.Config{
		ApiKeyValidity: 3 * time.Hour,
	}
	assert.Equal(expected, conf.ApiKey)
}

const defaultConfigName = "some-config"

func TestLoadConfiguration_LooksForTheRightFile(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	loadConf, actualConfigName := generateLoadConfFunc(nil)

	loadConf(defaultConfigName, defaultConf())

	assert.Equal(defaultConfigName, *actualConfigName)
}

func TestLoadConfiguration_WhenError_ExpectDefaultAndError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	loadConf, _ := generateLoadConfFunc(errDefault)

	conf, err := loadConf(defaultConfigName, defaultConf())

	assert.Equal(errDefault, err)
	assert.Equal(defaultConf(), conf)
}

func resetConfigurationParser() {
	loadConfFunc = config.LoadConfiguration[Configuration]
}

func generateLoadConfFunc(loadErr error) (loadConfFuncType, *string) {
	var outConfName string

	f := func(configName string, defaultConf Configuration) (Configuration, error) {
		outConfName = configName
		return defaultConf, loadErr
	}

	return f, &outConfName
}
