package controllers

import (
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/labstack/echo/v4"
)

type repositoriesAwareHttpHandler func(echo.Context, repositories.UserRepository) error

func generateEchoHandler(handler repositoriesAwareHttpHandler, repo repositories.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, repo)
	}
}
