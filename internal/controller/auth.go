package controller

import (
	"encoding/json"
	"net/http"

	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/middleware"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/rest"
	"github.com/labstack/echo/v4"
)

const aclHeaderKey = "X-Acl"
const userLimitHeaderKey = "X-User-Limit"

func AuthEndpoints(service service.AuthService) rest.Routes {
	var out rest.Routes

	authHandler := fromAuthServiceAwareHttpHandler(authUser, service)
	auth := rest.NewRoute(http.MethodGet, "/auth", authHandler)
	out = append(out, auth)

	return out
}

func authUser(c echo.Context, s service.AuthService) error {
	apiKey, err := middleware.TryGetApiKeyHeader(c.Request())
	if err != nil {
		if errors.IsErrorWithCode(err, middleware.ApiKeyNotFound) || errors.IsErrorWithCode(err, middleware.TooManyApiKeys) {
			return c.JSON(http.StatusBadRequest, "Api key not found")
		} else if errors.IsErrorWithCode(err, middleware.InvalidApiKeySyntax) {
			return c.JSON(http.StatusBadRequest, "Invalid api key syntax")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	out, err := s.Authenticate(c.Request().Context(), apiKey)
	if err != nil {
		if errors.IsErrorWithCode(err, service.UserNotAuthenticated) {
			return c.JSON(http.StatusForbidden, err)
		} else if errors.IsErrorWithCode(err, service.AuthenticationExpired) {
			return c.JSON(http.StatusForbidden, err)
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	aclJson, err := json.Marshal(out.Acls)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set(aclHeaderKey, string(aclJson))

	userLimitJson, err := json.Marshal(out.Limits)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set(userLimitHeaderKey, string(userLimitJson))

	return c.NoContent(http.StatusNoContent)
}
