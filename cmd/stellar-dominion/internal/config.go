package internal

import (
	"github.com/KnoblauchPilze/user-service/pkg/config"
)

func DefaultConf() config.Configuration {
	conf := config.DefaultConf()
	conf.Database.Name = "db_stellar_dominion"
	return conf
}
