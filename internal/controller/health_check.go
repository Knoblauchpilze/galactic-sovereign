package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/rest"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/labstack/echo/v4"
)

func HealthCheckEndpoints(pool db.ConnectionPool) rest.Routes {
	var out rest.Routes

	getHandler := fromDbAwareHttpHandler(healthcheck, pool)
	get := rest.NewRoute(http.MethodGet, "/healthcheck", getHandler)
	out = append(out, get)

	return out
}

func healthcheck(c echo.Context, pool db.ConnectionPool) error {
	err := pool.Ping(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, errors.WrapCode(err, HealthcheckFailed))
	}

	return c.JSON(http.StatusOK, "OK")
}
