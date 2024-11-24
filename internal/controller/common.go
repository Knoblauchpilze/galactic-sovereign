package controller

import (
	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/labstack/echo/v4"
)

type buildingActionServiceAwareHttpHandler func(echo.Context, service.BuildingActionService) error

func fromBuildingActionServiceAwareHttpHandler(handler buildingActionServiceAwareHttpHandler, service service.BuildingActionService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}

type dbAwareHttpHandler func(echo.Context, db.ConnectionPool) error

func fromDbAwareHttpHandler(handler dbAwareHttpHandler, pool db.ConnectionPool) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, pool)
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
