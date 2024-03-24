package middleware

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type keyType string

const requestIdKey keyType = "requestIdKey"
const logKey keyType = "loggerKey"

func ResponseEnvelope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := uuid.Must(uuid.NewRandom())
			log := logger.New(id.String())

			c.Set(string(requestIdKey), id)
			c.SetLogger(log)

			ctx := context.WithValue(c.Request().Context(), logKey, log)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			w := new(c.Response().Writer, id, log)
			c.Response().Writer = w

			return next(c)
		}
	}
}

func GetLoggerFromContext(ctx context.Context) echo.Logger {
	log, ok := ctx.Value(logKey).(echo.Logger)
	if ok && log != nil {
		return log
	}

	id, ok := ctx.Value(requestIdKey).(uuid.UUID)
	if ok {
		return logger.New(id.String())
	}

	return logger.New("")
}
