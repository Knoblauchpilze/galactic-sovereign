package internal

import (
	"github.com/KnoblauchPilze/user-service/internal/config"
)

func DefaultConf() config.Configuration {
	conf := config.DefaultConf()
	conf.Database.Name = "db_stellar_dominion"
	return conf
}
