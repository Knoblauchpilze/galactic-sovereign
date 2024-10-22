package internal

import (
	"github.com/KnoblauchPilze/galactic-sovereign/internal/config"
)

func DefaultConf() config.Configuration {
	conf := config.DefaultConf()
	conf.Database.Name = "db_galactic_sovereign"
	return conf
}
