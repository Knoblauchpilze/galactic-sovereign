package internal

import (
	"github.com/KnoblauchPilze/galactic-sovereign/internal/config"
)

func DefaultConf() config.Configuration {
	conf := config.DefaultConf()
	conf.Server.BasePath = "/v1/galactic-sovereign"
	conf.Database.Name = "db_galactic_sovereign"
	return conf
}
