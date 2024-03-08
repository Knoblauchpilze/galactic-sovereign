package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	logMiddleware := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_rfc3339_nano} ${method} ${uri} -> ${status}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           os.Stdout,
	})
	e.Use(logMiddleware)
	e.Use(middleware.Recover())

	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
