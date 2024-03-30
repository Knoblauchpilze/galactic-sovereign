package internal

import (
	"strings"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   rest.Config
	Database db.Config
}

func LoadConfiguration() (Configuration, error) {
	// https://github.com/spf13/viper#reading-config-files
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	viper.SetConfigName("users-dev")
	if err := viper.ReadInConfig(); err != nil {
		return defaultConf(), err
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("ENV")
	viper.AutomaticEnv()

	var out Configuration
	if err := viper.Unmarshal(&out); err != nil {
		return defaultConf(), err
	}

	return out, nil
}

func defaultConf() Configuration {
	return Configuration{
		Server: rest.Config{
			Endpoint: "/v1/users/",
			Port:     uint16(60000),
		},
		Database: db.Config{
			Host:                "localhost",
			Port:                5432,
			Name:                "database",
			User:                "user",
			Password:            "password",
			ConnectionsPoolSize: 1,
		},
	}
}
