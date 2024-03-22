package db

import (
	"github.com/jackc/pgx"
)

type Config struct {
	Host                string
	Port                uint16
	Name                string
	User                string
	Password            string
	ConnectionsPoolSize uint
}

func (c Config) toConnPoolConfig() pgx.ConnPoolConfig {
	return pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     c.Host,
			Port:     c.Port,
			Database: c.Name,
			User:     c.User,
			Password: c.Password,
		},
		// TODO: Also set the logger?
		// Logger            Logger
		// LogLevel          LogLevel

		MaxConnections: int(c.ConnectionsPoolSize),
		AcquireTimeout: 0,
	}
}
