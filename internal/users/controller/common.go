package controller

import (
	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/labstack/echo/v4"
)

type repositoriesAwareHttpHandler func(echo.Context, service.UserService) error

func fromRepositoriesAwareHttpHandler(handler repositoriesAwareHttpHandler, service service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}
