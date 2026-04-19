package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/config"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/logger"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/process"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	docs "github.com/Knoblauchpilze/galactic-sovereign/api"
	"github.com/Knoblauchpilze/galactic-sovereign/cmd/galactic-sovereign/internal"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/controller"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func determineConfigName() string {
	if len(os.Args) < 2 {
		return "galactic-sovereign-prod.yml"
	}

	return os.Args[1]
}

// @title			Galactic Sovereign API
// @version		1.0
// @description	REST API for the Galactic Sovereign backend service.
// @BasePath		/v1/galactic-sovereign
func main() {
	log := logger.New(os.Stdout)

	conf, err := config.Load(determineConfigName(), internal.DefaultConfig())
	if err != nil {
		log.Error("Failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	conn, err := db.New(context.Background(), conf.Database)
	if err != nil {
		log.Error("Failed to create db connection", slog.Any("error", err))
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
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}

	for _, route := range controller.PlayerEndpoints(playerService) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}

	for _, route := range controller.UniverseEndpoints(universeService) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}

	for _, route := range controller.BuildingActionEndpoints(buildingActionService, actionService, planetResourceService) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}

	for _, route := range controller.HealthCheckEndpoints(conn) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}

	docs.SwaggerInfo.BasePath = conf.Server.BasePath
	swaggerUi := rest.NewRawRoute(http.MethodGet, "/swagger/*", echoSwagger.WrapHandler)
	if err := s.AddRoute(swaggerUi); err != nil {
		log.Error("Failed to register route", slog.String("route", swaggerUi.Path()), slog.Any("error", err))
		os.Exit(1)
	}

	wait, err := process.StartWithSignalHandler(context.Background(), s)
	if err != nil {
		log.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}

	err = wait()
	if err != nil {
		log.Error("Error while serving", slog.Any("error", err))
		os.Exit(1)
	}
}
