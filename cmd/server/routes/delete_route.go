package routes

import "github.com/labstack/echo/v4"

func registerDeleteRoute(path string, handler echo.HandlerFunc, e *echo.Echo) {
	deletePath := concatenateEndpoints(path, ":id")
	e.Logger.Debugf("Adding route DELETE %s", deletePath)
	e.DELETE(deletePath, handler)
}
