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
	drivenadapters "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven"
	drivingadapters "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases"
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

	s := server.NewWithLogger(conf.Server, log)

	// TODO: Some of the routes need to be restored under the game.NewResourceRoute
	// wrapping to trigger the processing of actions (e.g. planets)
	registerUniversesRoutes(conn, s, log)
	registerPlayersRoutes(conn, s, log)
	registerPlanetsRoutes(conn, s, log)
	registerBuildingActionsRoutes(conn, s, log)
	registerHealthRoutes(conn, s, log)

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

func registerUniversesRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	repo := drivenadapters.NewUniverseRepository(conn)
	usecase := usecases.NewUniverseUseCase(repo)

	for _, route := range drivingadapters.UniverseEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}
}

func registerPlayersRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	playerRepo := drivenadapters.NewPlayerRepository(conn)
	universeRepo := drivenadapters.NewUniverseRepository(conn)
	planetRepo := drivenadapters.NewPlanetRepository(conn)
	usecase := usecases.NewPlayerUseCase(playerRepo, universeRepo, planetRepo)

	for _, route := range drivingadapters.PlayerEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}
}

func registerPlanetsRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	playerRepo := drivenadapters.NewPlayerRepository(conn)
	planetRepo := drivenadapters.NewPlanetRepository(conn)

	usecase := usecases.NewPlanetUseCase(playerRepo, planetRepo)

	for _, route := range drivingadapters.PlanetEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}
}

func registerBuildingActionsRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	actionRepo := drivenadapters.NewBuildingActionRepository(conn)
	planetRepo := drivenadapters.NewPlanetRepository(conn)
	buildingRepo := drivenadapters.NewBuildingRepository(conn)
	usecase := usecases.NewBuildingActionUseCase(actionRepo, planetRepo, buildingRepo)

	for _, route := range drivingadapters.BuildingActionEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}
}

func registerHealthRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	checker := drivenadapters.NewDatabaseChecker(conn)
	usecase := usecases.NewCheckHealthUseCase(checker)

	for _, route := range drivingadapters.HealthcheckEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}
}
