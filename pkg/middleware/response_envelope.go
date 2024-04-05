package middleware

import (
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ResponseEnvelope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqId := uuid.Must(uuid.NewRandom())

			log, ctx := logger.RegisterRequestLogger(reqId, c.Request().Context())
			c.SetLogger(log)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			w := new(c.Response().Writer, reqId, log)
			c.Response().Writer = w

			return next(c)
		}
	}
}
