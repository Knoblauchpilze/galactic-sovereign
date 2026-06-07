package drivingadapters

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/labstack/echo/v5"
)

func HealthcheckEndpoints(usecase drivingports.ForCheckingServiceHealth) rest.Routes {
	var out rest.Routes

	handler := generateHandler(healthcheck, usecase)
	get := rest.NewRoute(http.MethodGet, "/healthcheck", handler)
	out = append(out, get)

	return out
}

// healthcheck godoc
//
//	@Summary		Health check
//	@Description	Returns service health based on database connectivity.
//	@Tags			healthcheck
//	@Produce		json
//	@Success		200	{object}	rest.ResponseEnvelope[string]	"OK"
//	@Failure		503	{object}	rest.ResponseEnvelope[string]
//	@Router			/healthcheck [get]
func healthcheck(c *echo.Context, usecase drivingports.ForCheckingServiceHealth) error {
	healthy := usecase.Healthy(c.Request().Context())

	if !healthy {
		return c.JSON(http.StatusServiceUnavailable, "KO")
	}

	return c.JSON(http.StatusOK, "OK")
}
