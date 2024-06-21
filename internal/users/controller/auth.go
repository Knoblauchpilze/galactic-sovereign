package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/echo/v4"
)

func AuthEndpoints(service service.AuthService) rest.Routes {
	var out rest.Routes

	authHandler := fromAuthServiceAwareHttpHandler(authUser, service)
	auth := rest.NewRoute(http.MethodGet, false, "/users/auth", authHandler)
	out = append(out, auth)

	return out
}

func authUser(c echo.Context, as service.AuthService) error {
	apiKey, err := middleware.TryGetApiKeyHeader(c.Request())
	if err != nil {
		if errors.IsErrorWithCode(err, middleware.ApiKeyNotFound) || errors.IsErrorWithCode(err, middleware.TooManyApiKeys) {
			return c.JSON(http.StatusBadRequest, "Api key not found")
		} else if errors.IsErrorWithCode(err, middleware.InvalidApiKeySyntax) {
			return c.JSON(http.StatusBadRequest, "Invalid api key syntax")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	out, err := as.Authenticate(c.Request().Context(), apiKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}
