package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UserRoutes() Route {
	return Route{
		Path:        "users",
		GetRoute:    getUser,
		PostRoute:   createUser,
		DeleteRoute: deleteUser,
	}
}

func getUser(c echo.Context) error {
	panic(fmt.Errorf("not implemented"))
}

func createUser(c echo.Context) error {
	panic(fmt.Errorf("not implemented"))
}

func deleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "{\"message\":\"deleted\"}")
}
