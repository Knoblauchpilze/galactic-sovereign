package internal

import (
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
)

type Configuration struct {
	Server   server.Config
	Port     uint16
	Database postgresql.Config
}

func DefaultConfig() Configuration {
	const defaultDatabaseName = "db_galactic_sovereign"
	const defaultDatabaseUser = "galactic_sovereign_manager"

	return Configuration{
		Server: server.Config{
			BasePath:        "/v1/galactic-sovereign",
			ShutdownTimeout: 5 * time.Second,
		},
		Port: uint16(80),
		Database: postgresql.NewConfigForDockerContainer(
			defaultDatabaseName,
			defaultDatabaseUser,
			"comes-from-the-environment",
		),
	}
}
