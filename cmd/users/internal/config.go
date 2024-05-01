package internal

import (
	"strings"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   rest.Config
	Database db.Config
}

type configurationParser interface {
	SetConfigType(extension string)
	AddConfigPath(path string)

	SetConfigName(fileName string)
	ReadInConfig() error

	SetEnvKeyReplacer(replacer *strings.Replacer)
	SetEnvPrefix(envPrefix string)
	AutomaticEnv()

	Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error
}

var configurator configurationParser = viper.New()

func LoadConfiguration(configName string) (Configuration, error) {
	// https://github.com/spf13/viper#reading-config-files
	configurator.SetConfigType("yaml")
	configurator.AddConfigPath("configs")

	// https://stackoverflow.com/questions/61585304/issues-with-overriding-config-using-env-variables-in-viper
	configurator.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configurator.SetEnvPrefix("ENV")
	configurator.AutomaticEnv()

	configurator.SetConfigName(configName)
	if err := configurator.ReadInConfig(); err != nil {
		return defaultConf(), err
	}

	out := defaultConf()
	if err := configurator.Unmarshal(&out); err != nil {
		return defaultConf(), err
	}

	return out, nil
}

func defaultConf() Configuration {
	return Configuration{
		Server: rest.Config{
			BasePath:  "/v1/users/",
			Port:      uint16(80),
			RateLimit: 10,
		},
		Database: db.Config{
			Host:                "172.17.0.1",
			Port:                5432,
			Name:                "db_user_service",
			ConnectionsPoolSize: 1,
			LogLevel:            log.DEBUG,
		},
	}
}
