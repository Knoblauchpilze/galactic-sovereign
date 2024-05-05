package rest

import "github.com/labstack/echo/v4"

type echoServerImpl struct {
	e *echo.Echo
}

func (esi *echoServerImpl) Use(m ...echo.MiddlewareFunc) {
	esi.e.Use(m...)
}

func (esi *echoServerImpl) Group(prefix string, m ...echo.MiddlewareFunc) echoRouter {
	return esi.e.Group(prefix, m...)
}

func (esi *echoServerImpl) Start(address string) error {
	return esi.e.Start(address)
}
