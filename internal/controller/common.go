package controller

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	"github.com/labstack/echo/v5"
)

type buildingActionServiceAwareHttpHandler func(*echo.Context, service.BuildingActionService) error

func fromBuildingActionServiceAwareHttpHandler(handler buildingActionServiceAwareHttpHandler, service service.BuildingActionService) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return handler(c, service)
	}
}

type dbAwareHttpHandler func(*echo.Context, db.Connection) error

func fromDbAwareHttpHandler(handler dbAwareHttpHandler, conn db.Connection) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return handler(c, conn)
	}
}
