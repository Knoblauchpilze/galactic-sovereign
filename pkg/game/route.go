package game

import (
	"strings"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/labstack/echo/v4"
)

func NewRoute(method string,
	path string,
	handler echo.HandlerFunc,
	actionService ActionService,
	planetResourceService PlanetResourceService) rest.Route {
	wrapped := GameUpdateWatcher(actionService, planetResourceService, handler)
	return rest.NewRoute(method, path, wrapped)
}

func NewResourceRoute(method string,
	path string,
	handler echo.HandlerFunc,
	actionService ActionService,
	planetResourceService PlanetResourceService) rest.Route {
	wrapped := GameUpdateWatcher(actionService, planetResourceService, handler)

	if !strings.Contains(path, "/:id") {
		path = rest.ConcatenateEndpoints(path, "/:id")
	}

	return rest.NewRoute(method, path, wrapped)
}
