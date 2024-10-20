package middleware

import (
	"net/http"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const apiKeyHeaderKey = "X-Api-Key"

const (
	ApiKeyNotFound      errors.ErrorCode = 200
	TooManyApiKeys      errors.ErrorCode = 201
	InvalidApiKeySyntax errors.ErrorCode = 202
)

func ApiKey(apiKeyRepository repositories.ApiKeyRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKeyValue, err := TryGetApiKeyHeader(c.Request())
			if err != nil {
				c.Logger().Errorf("Failed to fetch key: %v", err)

				if errors.IsErrorWithCode(err, InvalidApiKeySyntax) {
					return c.JSON(http.StatusBadRequest, "API key has wrong format")
				}

				return c.JSON(http.StatusBadRequest, "API key not found")
			}

			apiKey, err := apiKeyRepository.GetForKey(c.Request().Context(), apiKeyValue)
			if err != nil {
				if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
					return c.JSON(http.StatusUnauthorized, "Invalid API key")
				}

				c.Logger().Warnf("Failed to fetch key %v: %v", apiKeyValue, err)
				return c.JSON(http.StatusInternalServerError, "Failed to verify API key")
			}

			if time.Now().After(apiKey.ValidUntil) {
				c.Logger().Errorf("API Key %v expired since %v", apiKey.Id, apiKey.ValidUntil)
				return c.JSON(http.StatusUnauthorized, "API key expired")
			}

			return next(c)
		}
	}
}

func TryGetApiKeyHeader(req *http.Request) (apiKey uuid.UUID, err error) {
	apiKeys, ok := req.Header[apiKeyHeaderKey]
	if !ok {
		err = errors.NewCode(ApiKeyNotFound)
		return
	}
	if len(apiKeys) != 1 {
		err = errors.NewCode(TooManyApiKeys)
		return
	}

	apiKey, err = uuid.Parse(apiKeys[0])
	if err != nil {
		err = errors.NewCode(InvalidApiKeySyntax)
	}

	return
}
