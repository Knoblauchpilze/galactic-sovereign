package game

import (
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/rest"
	"github.com/labstack/echo/v4"
)

func NewRoute(method string,
	path string,
	handler echo.HandlerFunc,
	actionService ActionService,
	planetResourceService PlanetResourceService) rest.Route {
	wrapped := GameUpdateWatcher(actionService, planetResourceService, handler)
	return rest.NewRoute(method, false, path, wrapped)
}

func NewResourceRoute(method string,
	path string,
	handler echo.HandlerFunc,
	actionService ActionService,
	planetResourceService PlanetResourceService) rest.Route {
	wrapped := GameUpdateWatcher(actionService, planetResourceService, handler)
	return rest.NewResourceRoute(method, false, path, wrapped)
}
