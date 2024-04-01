package controller

import (
	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/labstack/echo/v4"
)

type repositoriesAwareHttpHandler func(echo.Context, repositories.UserRepository, service.UserService) error

func generateEchoHandler(handler repositoriesAwareHttpHandler, repo repositories.UserRepository, users service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, repo, users)
	}
}
