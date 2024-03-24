package middleware

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const requestIdKey = "requestIdKey"

type loggerKey string

const logKey loggerKey = "loggerKey"

func ResponseEnvelope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := uuid.Must(uuid.NewRandom())
			log := logger.New(id.String())

			c.Set(requestIdKey, id)
			c.SetLogger(log)

			ctx := context.WithValue(c.Request().Context(), logKey, log)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			w := new(c.Response().Writer)
			c.Response().Writer = w

			return next(c)
		}
	}
}

func GetLoggerFromContext(ctx context.Context) echo.Logger {
	log := ctx.Value(logKey).(echo.Logger)
	if log != nil {
		return log
	}

	id := ctx.Value(requestIdKey).(uuid.UUID)

	return logger.New(id.String())
}
