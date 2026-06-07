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
	_ "github.com/Knoblauchpilze/galactic-sovereign/api"
	"github.com/Knoblauchpilze/galactic-sovereign/cmd/galactic-sovereign/internal"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/controller"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	drivenadapters "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven"
	drivingadapters "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases"
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
// @servers.url /v1
// @servers.description Base path for the galactic-sovereign API
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
		PlanetBuilding:                   repositories.NewPlanetBuildingRepository(),
		PlanetResource:                   repositories.NewPlanetResourceRepository(),
		PlanetResourceProduction:         repositories.NewPlanetResourceProductionRepository(),
		PlanetResourceStorage:            repositories.NewPlanetResourceStorageRepository(),
		Resource:                         repositories.NewResourceRepository(),
	}

	buildingActionService := service.NewBuildingActionService(conn, repos)

	actionService := service.NewActionService(conn, repos)
	planetResourceService := service.NewPlanetResourceService(conn, repos)

	s := server.NewWithLogger(conf.Server, log)

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

	// New logic using DDD
	registerUniverseRoutes(conn, s, log)
	registerPlayerRoutes(conn, s, log)
	registerPlanetsRoutes(conn, s, log)
	// End new logic using DDD

	swaggerUi := rest.NewRawRoute(http.MethodGet, "/swagger/*", echoSwagger.WrapHandlerV3)
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

func registerUniverseRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	repo := drivenadapters.NewUniverseRepository(conn)
	usecase := usecases.NewUniverseUseCase(repo)

	for _, route := range drivingadapters.UniverseEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}
}

func registerPlayerRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	repo := drivenadapters.NewPlayerRepository(conn)
	usecase := usecases.NewPlayerUseCase(repo)

	for _, route := range drivingadapters.PlayerEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}
}

func registerPlanetsRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	repo := drivenadapters.NewPlanetRepository(conn)
	usecase := usecases.NewPlanetUseCase(repo)

	for _, route := range drivingadapters.PlanetEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}
}
