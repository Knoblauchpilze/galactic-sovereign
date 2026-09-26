package internal

import (
	"context"
	"log/slog"
	"net"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	drivenadapters "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	drivingadapters "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases"
)

type HttpServer interface {
	AddRoute(route *rest.Route) error
	Bind(port uint16) (net.Listener, error)
	Serve(ctx context.Context, listener net.Listener) error
}

func CreateGameServer(
	conf server.Config,
	conn database.Connection,
	log *slog.Logger,
) HttpServer {
	s := server.NewHttpServerWithLogger(conf, log)

	registerUniversesRoutes(conn, s, log)
	registerPlayersRoutes(conn, s, log)
	registerPlanetsRoutes(conn, s, log)
	registerBuildingActionsRoutes(conn, s, log)
	registerShipsRoutes(conn, s, log)
	registerHealthRoutes(conn, s, log)

	return s
}

func registerUniversesRoutes(
	conn database.Connection,
	s HttpServer,
	log *slog.Logger,
) {
	universeRepo := drivenadapters.NewUniverseRepository(conn)
	solarSystemRepo := drivenadapters.NewSolarSystemRepository(conn)

	universeUsecase := usecases.NewUniverseUseCase(universeRepo)
	listSolarSystemUsecase := usecases.NewFetchSolarSystemUseCase(solarSystemRepo)

	for _, route := range drivingadapters.UniverseEndpoints(universeUsecase, listSolarSystemUsecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
		}
	}
}

func registerPlayersRoutes(
	conn database.Connection,
	s HttpServer,
	log *slog.Logger,
) {
	playerRepo := drivenadapters.NewPlayerRepository(conn)
	universeRepo := drivenadapters.NewUniverseRepository(conn)
	usecase := usecases.NewPlayerUseCase(playerRepo, universeRepo)

	for _, route := range drivingadapters.PlayerEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
		}
	}
}

func registerPlanetsRoutes(
	conn database.Connection,
	s HttpServer,
	log *slog.Logger,
) {
	planetRepo := drivenadapters.NewPlanetRepository(conn)
	planetMutator := drivenadapters.NewPlanetMutator(conn)
	clock := drivenadapters.NewTimeAdapter()

	usecase := usecases.NewPlanetUseCase(planetRepo, planetMutator, clock)

	for _, route := range drivingadapters.PlanetEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
		}
	}
}

func registerBuildingActionsRoutes(
	conn database.Connection,
	s HttpServer,
	log *slog.Logger,
) {
	buildingRepo := drivenadapters.NewBuildingRepository(conn)
	planetMutator := drivenadapters.NewPlanetMutator(conn)
	clock := drivenadapters.NewTimeAdapter()

	createUseCase := usecases.NewCreateBuildingActionUseCase(buildingRepo, planetMutator, clock)
	deleteUsecase := usecases.NewDeleteBuildingActionUseCase(planetMutator, clock)

	for _, route := range drivingadapters.BuildingActionEndpoints(createUseCase, deleteUsecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
		}
	}
}

func registerShipsRoutes(
	conn database.Connection,
	s HttpServer,
	log *slog.Logger,
) {
	shipRepo := drivenadapters.NewShipRepository(conn)
	planetMutator := drivenadapters.NewPlanetMutator(conn)
	clock := drivenadapters.NewTimeAdapter()

	createUseCase := usecases.NewCreateShipActionUseCase(shipRepo, planetMutator, clock)

	for _, route := range drivingadapters.ShipActionEndpoints(createUseCase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
		}
	}
}

func registerHealthRoutes(
	conn database.Connection,
	s HttpServer,
	log *slog.Logger,
) {
	checker := drivenadapters.NewDatabaseChecker(conn)
	usecase := usecases.NewCheckHealthUseCase(checker)

	for _, route := range drivingadapters.HealthcheckEndpoints(usecase) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
		}
	}
}
