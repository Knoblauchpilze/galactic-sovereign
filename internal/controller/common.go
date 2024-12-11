package controller

import (
	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/labstack/echo/v4"
)

type buildingActionServiceAwareHttpHandler func(echo.Context, service.BuildingActionService) error

func fromBuildingActionServiceAwareHttpHandler(handler buildingActionServiceAwareHttpHandler, service service.BuildingActionService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}

type dbAwareHttpHandler func(echo.Context, db.Connection) error

func fromDbAwareHttpHandler(handler dbAwareHttpHandler, conn db.Connection) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, conn)
	}
}

type planetServiceAwareHttpHandler func(echo.Context, service.PlanetService) error

func fromPlanetServiceAwareHttpHandler(handler planetServiceAwareHttpHandler, service service.PlanetService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}

type playerServiceAwareHttpHandler func(echo.Context, service.PlayerService) error

func fromPlayerServiceAwareHttpHandler(handler playerServiceAwareHttpHandler, service service.PlayerService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}

type universeServiceAwareHttpHandler func(echo.Context, service.UniverseService) error

func fromUniverseServiceAwareHttpHandler(handler universeServiceAwareHttpHandler, service service.UniverseService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}
