package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KnoblauchPilze/user-service/cmd/users/internal"
	"github.com/KnoblauchPilze/user-service/internal/users/controller"
	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
)

func main() {
	conf, err := internal.LoadConfiguration()
	if err != nil {
		logger.Errorf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	pool := db.NewConnectionPool(conf.Database)
	if err := pool.Connect(); err != nil {
		logger.Errorf("Failed to connect to the database: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	installCleanup(pool)

	userRepo := repositories.NewUserRepository(pool)
	apiKeyRepo := repositories.NewApiKeyRepository(pool)
	userService := service.NewUserService(pool, userRepo, apiKeyRepo)

	s := rest.NewServer(conf.Server)

	for _, route := range controller.UserEndpoints(userService) {
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

func installCleanup(conn db.ConnectionPool) {
	// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i
	interruptChannel := make(chan os.Signal, 2)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChannel

		conn.Close()
		os.Exit(1)
	}()
}
