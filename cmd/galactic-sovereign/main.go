package main

import (
	"context"
	"os"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/config"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/logger"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	"github.com/Knoblauchpilze/galactic-sovereign/cmd/galactic-sovereign/internal"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/controller"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
)

func determineConfigName() string {
	if len(os.Args) < 2 {
		return "galactic-sovereign-prod.yml"
	}

	return os.Args[1]
}

func main() {
	log := logger.New(logger.NewPrettyWriter(os.Stdout))

	conf, err := config.Load(determineConfigName(), internal.DefaultConfig())
	if err != nil {
		log.Errorf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	conn, err := db.New(context.Background(), conf.Database)
	if err != nil {
		log.Errorf("Failed to create db connection: %v", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	repos := repositories.Repositories{
		Building:                         repositories.NewBuildingRepository(),
		BuildingAction:                   repositories.NewBuildingActionRepository(),
		BuildingActionCost:               repositories.NewBuildingActionCostRepository(),
		BuildingActionResourceProduction: repositories.NewBuildingActionResourceProductionRepository(),
		BuildingActionResourceStorage:    repositories.NewBuildingActionResourceStorageRepository(),
		BuildingCost:                     repositories.NewBuildingCostRepository(),
		BuildingResourceProduction:       repositories.NewBuildingResourceProductionRepository(),
		BuildingResourceStorage:          repositories.NewBuildingResourceStorageRepository(),
		Planet:                           repositories.NewPlanetRepository(conn),
		PlanetBuilding:                   repositories.NewPlanetBuildingRepository(),
		PlanetResource:                   repositories.NewPlanetResourceRepository(),
		PlanetResourceProduction:         repositories.NewPlanetResourceProductionRepository(),
		PlanetResourceStorage:            repositories.NewPlanetResourceStorageRepository(),
		Player:                           repositories.NewPlayerRepository(conn),
		Resource:                         repositories.NewResourceRepository(),
		Universe:                         repositories.NewUniverseRepository(conn),
	}

	planetService := service.NewPlanetService(conn, repos)
	playerService := service.NewPlayerService(conn, repos)
	universeService := service.NewUniverseService(conn, repos)
	buildingActionService := service.NewBuildingActionService(conn, repos)

	actionService := service.NewActionService(conn, repos)
	planetResourceService := service.NewPlanetResourceService(conn, repos)

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

	for _, route := range controller.HealthCheckEndpoints(conn) {
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
