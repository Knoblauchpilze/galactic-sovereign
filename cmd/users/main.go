package main

import (
	"context"
	"fmt"
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

func determineConfigName() string {
	if len(os.Args) < 2 {
		return "users-prod.yml"
	}

	return os.Args[1]
}

func main() {
	conf, err := internal.LoadConfiguration(determineConfigName())
	if err != nil {
		logger.Errorf("Failed to load configuration: %v", err)
		os.Exit(1)
	}
	if conf.Database.User == "" || conf.Database.Password == "" {
		logger.Errorf("Please provide a user and password to connect to the database")
		os.Exit(1)
	}

	fmt.Printf("conf: %+v\n", conf)

	pool := db.NewConnectionPool(conf.Database)
	if err := pool.Connect(context.Background()); err != nil {
		logger.Errorf("Failed to connect to the database: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Errorf("Failed to ping the database: %v", err)
		os.Exit(1)
	}

	userRepo := repositories.NewUserRepository(pool)
	apiKeyRepo := repositories.NewApiKeyRepository(pool)
	userService := service.NewUserService(conf.ApiKey, pool, userRepo, apiKeyRepo)

	s := rest.NewServer(conf.Server, apiKeyRepo)

	for _, route := range controller.UserEndpoints(userService) {
		if err := s.Register(route); err != nil {
			logger.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	for _, route := range controller.HealthCheckEndpoints(pool) {
		if err := s.Register(route); err != nil {
			logger.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	s.Start()

	installCleanup(pool, s)

	if err := s.Wait(); err != nil {
		logger.Errorf("Error while servier was running: %v", err)
		os.Exit(1)
	}
}

func installCleanup(conn db.ConnectionPool, s rest.Server) {
	// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i
	interruptChannel := make(chan os.Signal, 2)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChannel

		s.Stop()
		conn.Close()
		os.Exit(1)
	}()
}
