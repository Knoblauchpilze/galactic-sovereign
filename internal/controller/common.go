package controller

import (
	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/labstack/echo/v4"
)

type userServiceAwareHttpHandler func(echo.Context, service.UserService) error

func fromUserServiceAwareHttpHandler(handler userServiceAwareHttpHandler, service service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}

type dbAwareHttpHandler func(echo.Context, db.ConnectionPool) error

func fromDbAwareHttpHandler(handler dbAwareHttpHandler, pool db.ConnectionPool) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, pool)
	}
}

type authServiceAwareHttpHandler func(echo.Context, service.AuthService) error

func fromAuthServiceAwareHttpHandler(handler authServiceAwareHttpHandler, service service.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}
