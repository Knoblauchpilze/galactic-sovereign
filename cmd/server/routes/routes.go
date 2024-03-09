package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

func setupRoutes(e *echo.Echo) {
	e.GET("/hello", makeHello)
	e.GET("/panic", makePanic)
	e.GET("/302", make302)
	e.GET("/418", make418)
	e.GET("/503", make503)
}

func makeHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func makePanic(c echo.Context) error {
	panic("a panic")
}

func make302(c echo.Context) error {
	c.Response().Status = http.StatusFound
	return nil
}

func make418(c echo.Context) error {
	c.Response().Status = http.StatusTeapot
	return nil
}

func make503(c echo.Context) error {
	c.Response().Status = http.StatusServiceUnavailable
	return nil
}
