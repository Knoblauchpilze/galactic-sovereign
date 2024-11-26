package main

import (
	"context"
	"os"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/logger"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/server"
	"github.com/KnoblauchPilze/galactic-sovereign/cmd/galactic-sovereign/internal"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/config"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/controller"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
)

func determineConfigName() string {
	if len(os.Args) < 2 {
		return "galactic-sovereign-prod.yml"
	}

	return os.Args[1]
}

func main() {
	log := logger.New(logger.NewPrettyWriter(os.Stdout))

	conf, err := config.LoadConfiguration(determineConfigName(), internal.DefaultConf())
	if err != nil {
		log.Errorf("Failed to load configuration: %v", err)
		os.Exit(1)
	}
	if conf.Database.User == "" || conf.Database.Password == "" {
		log.Errorf("Please provide a user and password to connect to the database")
		os.Exit(1)
	}

	pool := db.NewConnectionPool(conf.Database, log)
	if err := pool.Connect(context.Background()); err != nil {
		log.Errorf("Failed to connect to the database: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Errorf("Failed to ping the database: %v", err)
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

	s := server.NewWithLogger(conf.Server, log)

	for _, route := range controller.PlanetEndpoints(planetService, actionService, planetResourceService) {
		if err := s.AddRoute(route); err != nil {
			log.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	for _, route := range controller.PlayerEndpoints(playerService) {
		if err := s.AddRoute(route); err != nil {
			log.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	for _, route := range controller.UniverseEndpoints(universeService) {
		if err := s.AddRoute(route); err != nil {
			log.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	for _, route := range controller.BuildingActionEndpoints(buildingActionService, actionService, planetResourceService) {
		if err := s.AddRoute(route); err != nil {
			log.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	for _, route := range controller.HealthCheckEndpoints(pool) {
		if err := s.AddRoute(route); err != nil {
			log.Errorf("Failed to register route: %v", err)
			os.Exit(1)
		}
	}

	err = s.Start(context.Background())
	if err != nil {
		log.Errorf("Error while serving: %v", err)
		os.Exit(1)
	}
}
