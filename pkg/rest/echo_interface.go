package rest

import (
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/labstack/echo/v4"
)

type echoRouter interface {
	GET(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	POST(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	DELETE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	PATCH(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route

	Use(...echo.MiddlewareFunc)
}

type echoServer interface {
	Use(...echo.MiddlewareFunc)

	Group(string, ...echo.MiddlewareFunc) echoRouter

	Start(string) error
}

func createEchoServer() echoServer {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = logger.New("server")

	return &echoServerImpl{e: e}
}
