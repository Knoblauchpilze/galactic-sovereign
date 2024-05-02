package rest

import "github.com/labstack/echo/v4"

type echoServerImpl struct {
	e *echo.Echo
}

func (esi *echoServerImpl) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return esi.e.GET(path, h, m...)
}

func (esi *echoServerImpl) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return esi.e.POST(path, h, m...)
}

func (esi *echoServerImpl) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return esi.e.DELETE(path, h, m...)
}

func (esi *echoServerImpl) PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return esi.e.PATCH(path, h, m...)
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
