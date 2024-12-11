package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/rest"
	"github.com/labstack/echo/v4"
)

func HealthCheckEndpoints(conn db.Connection) rest.Routes {
	var out rest.Routes

	getHandler := fromDbAwareHttpHandler(healthcheck, conn)
	get := rest.NewRoute(http.MethodGet, "/healthcheck", getHandler)
	out = append(out, get)

	return out
}

func healthcheck(c echo.Context, conn db.Connection) error {
	err := conn.Ping(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, errors.WrapCode(err, HealthcheckFailed))
	}

	return c.JSON(http.StatusOK, "OK")
}
