package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/labstack/echo/v5"
)

func HealthCheckEndpoints(conn db.Connection) rest.Routes {
	var out rest.Routes

	getHandler := fromDbAwareHttpHandler(healthcheck, conn)
	get := rest.NewRoute(http.MethodGet, "/healthcheck", getHandler)
	out = append(out, get)

	return out
}

// healthcheck godoc
//
//	@Summary		Health check
//	@Description	Returns service health based on database connectivity.
//	@Tags			healthcheck
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Failure		503	{object}	rest.ResponseEnvelope[string]
//	@Router			/healthcheck [get]
func healthcheck(c *echo.Context, conn db.Connection) error {
	err := conn.Ping(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, errors.WrapCode(err, HealthcheckFailed))
	}

	return c.JSON(http.StatusOK, "OK")
}
