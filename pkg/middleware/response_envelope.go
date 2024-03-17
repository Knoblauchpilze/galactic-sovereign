package middleware

import (
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const RequestIdKey = "requestIdKey"

func ResponseEnvelope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := uuid.Must(uuid.NewRandom())
			log := logger.New(id.String())

			c.Set(RequestIdKey, id)
			c.SetLogger(log)

			w := new(c.Response().Writer)
			c.Response().Writer = w

			return next(c)
		}
	}
}
