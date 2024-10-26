package internal

import (
	"github.com/KnoblauchPilze/galactic-sovereign/internal/config"
)

func DefaultConf() config.Configuration {
	conf := config.DefaultConf()
	conf.Server.Prefix = "/users"
	conf.Database.Name = "db_user_service"
	return conf
}
