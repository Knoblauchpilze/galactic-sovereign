package internal

import (
	"log/slog"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	drivenadapters "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven"
	drivingadapters "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases"
)

func CreateGameServer(conf server.Config, conn db.Connection, log *slog.Logger) server.Server {
	s := server.NewWithLogger(conf, log)

	// TODO: Some of the routes need to be restored under the game.NewResourceRoute
	// wrapping to trigger the processing of actions (e.g. planets)
	registerUniversesRoutes(conn, s, log)
	registerPlayersRoutes(conn, s, log)
	registerPlanetsRoutes(conn, s, log)
	registerBuildingActionsRoutes(conn, s, log)
	registerHealthRoutes(conn, s, log)

	return s
}

func registerUniversesRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	repo := drivenadapters.NewUniverseRepository(conn)
	usecase := usecases.NewUniverseUseCase(repo)

	for _, route := range drivingadapters.UniverseEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
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
		}
	}
}

func registerPlanetsRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	playerRepo := drivenadapters.NewPlayerRepository(conn)
	universeRepo := drivenadapters.NewUniverseRepository(conn)
	planetRepo := drivenadapters.NewPlanetRepository(conn)
	planetMutator := drivenadapters.NewPlanetMutator(conn)

	createUsecase := usecases.NewCreatePlanetUseCase(playerRepo, universeRepo, planetRepo)
	usecase := usecases.NewPlanetUseCase(planetRepo, planetMutator)

	for _, route := range drivingadapters.PlanetEndpoints(createUsecase, usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
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
		}
	}
}

func registerHealthRoutes(conn db.Connection, s server.Server, log *slog.Logger) {
	checker := drivenadapters.NewDatabaseChecker(conn)
	usecase := usecases.NewCheckHealthUseCase(checker)

	for _, route := range drivingadapters.HealthcheckEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
		}
	}
}
