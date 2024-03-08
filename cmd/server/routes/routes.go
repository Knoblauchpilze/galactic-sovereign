package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

func setupRoutes(e *echo.Echo) {
	e.GET("/", hello)
}

// Handler
func hello(c echo.Context) error {
	c.Logger().Debugf("ok")
	c.Logger().Infof("ok")
	c.Logger().Warnf("ok")
	c.Logger().Errorf("ok")
	c.Logger().Printf("ok")
	return c.String(http.StatusOK, "Hello, World!\n")
}
