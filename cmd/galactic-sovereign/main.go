package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/KnoblauchPilze/galactic-sovereign/cmd/galactic-sovereign/internal"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/config"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/controller"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/logger"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/rest"
)

func determineConfigName() string {
	if len(os.Args) < 2 {
		return "galactic-sovereign-prod.yml"
	}

	return os.Args[1]
}

func main() {
	conf, err := config.LoadConfiguration(determineConfigName(), internal.DefaultConf())
	if err != nil {
		logger.Errorf("Failed to load configuration: %v", err)
		os.Exit(1)
	}
	if conf.Database.User == "" || conf.Database.Password == "" {
		logger.Errorf("Please provide a user and password to connect to the database")
		os.Exit(1)
	}

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

	repos := repositories.Repositories{
		Building:                         repositories.NewBuildingRepository(),
		BuildingAction:                   repositories.NewBuildingActionRepository(),
		BuildingActionCost:               repositories.NewBuildingActionCostRepository(),
		BuildingActionResourceProduction: repositories.NewBuildingActionResourceProductionRepository(),
		BuildingCost:                     repositories.NewBuildingCostRepository(),
		BuildingResourceProduction:       repositories.NewBuildingResourceProductionRepository(),
		Planet:                           repositories.NewPlanetRepository(pool),
		PlanetBuilding:                   repositories.NewPlanetBuildingRepository(),
		PlanetResource:                   repositories.NewPlanetResourceRepository(),
		PlanetResourceProduction:         repositories.NewPlanetResourceProductionRepository(),
		PlanetResourceStorage:            repositories.NewPlanetResourceStorageRepository(),
		Player:                           repositories.NewPlayerRepository(pool),
		Resource:                         repositories.NewResourceRepository(pool),
		Universe:                         repositories.NewUniverseRepository(pool),
	}

	planetService := service.NewPlanetService(pool, repos)
	playerService := service.NewPlayerService(pool, repos)
	universeService := service.NewUniverseService(pool, repos)
	buildingActionService := service.NewBuildingActionService(pool, repos)

	actionService := service.NewActionService(pool, repos)
	planetResourceService := service.NewPlanetResourceService(pool, repos)

	s := rest.NewServer(conf.Server)

	for _, route := range controller.PlanetEndpoints(planetService, actionService, planetResourceService) {
		if err := s.Register(route); err != nil {
			logger.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	for _, route := range controller.PlayerEndpoints(playerService) {
		if err := s.Register(route); err != nil {
			logger.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	for _, route := range controller.UniverseEndpoints(universeService) {
		if err := s.Register(route); err != nil {
			logger.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	for _, route := range controller.BuildingActionEndpoints(buildingActionService, actionService, planetResourceService) {
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
		logger.Errorf("Error while server was running: %v", err)
		os.Exit(1)
	}
}

func installCleanup(conn db.ConnectionPool, s rest.Server) {
	// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i
	interruptChannel := make(chan os.Signal, 2)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChannel

		if err := s.Stop(); err != nil {
			logger.Errorf("Error while stopping server: %v", err)
		}

		conn.Close()
		os.Exit(1)
	}()
}
