package internal

import (
	"time"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/config"
)

type Configuration struct {
	config.Configuration
	ApiKey service.Config
}

type loadConfFuncType func(string, Configuration) (Configuration, error)

var loadConfFunc = config.LoadConfiguration[Configuration]

func LoadConfiguration(configName string) (Configuration, error) {
	return loadConfFunc(configName, defaultConf())
}

func defaultConf() Configuration {
	conf := Configuration{
		Configuration: config.DefaultConf(),
		ApiKey: service.Config{
			ApiKeyValidity: time.Duration(3 * time.Hour),
		},
	}

	conf.Database.Name = "db_user_service"

	return conf
}
