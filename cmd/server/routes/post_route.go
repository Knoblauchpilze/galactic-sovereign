package routes

import "github.com/labstack/echo/v4"

func registerPostRoute(path string, handler echo.HandlerFunc, e *echo.Echo) {
	e.Logger.Debugf("Adding route POST %s", path)
	e.POST(path, handler)
}
