package game

import (
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/echo/v4"
)

func NewRoute(method string, path string, handler echo.HandlerFunc, actions ActionService) rest.Route {
	wrapped := GameUpdateWatcher(actions, handler)
	return rest.NewRoute(method, false, path, wrapped)
}

func NewResourceRoute(method string, path string, handler echo.HandlerFunc, actions ActionService) rest.Route {
	wrapped := GameUpdateWatcher(actions, handler)
	return rest.NewResourceRoute(method, false, path, wrapped)
}
