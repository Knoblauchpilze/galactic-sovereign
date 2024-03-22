package main

import (
	"os"

	"github.com/KnoblauchPilze/user-service/cmd/server/config"
	"github.com/KnoblauchPilze/user-service/cmd/server/routes"
	"github.com/KnoblauchPilze/user-service/cmd/server/server"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/logger"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		logger.Errorf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	conn := db.NewConnection(conf.Database)
	if err := conn.Connect(); err != nil {
		logger.Errorf("Failed to connect to the database: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	s := server.New(conf.Server)
	s.Register(routes.UserRoutes(conn))

	s.Start()
}
