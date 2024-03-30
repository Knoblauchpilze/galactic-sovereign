package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KnoblauchPilze/user-service/cmd/server/config"
	"github.com/KnoblauchPilze/user-service/cmd/server/controllers"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
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

	installCleanup(conn)

	s := rest.NewServer(conf.Server)

	for _, route := range controllers.UserEndpoints(conn) {
		if err := s.Register(route); err != nil {
			logger.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	if err := s.Start(); err != nil {
		logger.Errorf("Error while servier was running: %v", err)
		os.Exit(1)
	}
}

func installCleanup(conn db.Connection) {
	// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i
	interruptChannel := make(chan os.Signal, 2)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChannel

		conn.Close()
		os.Exit(1)
	}()
}
