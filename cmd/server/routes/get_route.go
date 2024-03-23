package routes

import "github.com/labstack/echo/v4"

func registerGetRoute(path string, handler echo.HandlerFunc, e *echo.Echo) {
	getPath := concatenateEndpoints(path, ":id")
	e.Logger.Debugf("Adding route GET %s", getPath)
	e.GET(getPath, handler)
}
