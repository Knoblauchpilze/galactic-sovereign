package routes

import "github.com/labstack/echo/v4"

type Route interface {
	Method() string
	Register(path string, e *echo.Echo)
}

type Routes []Route

type routeImpl struct {
	method  string
	path    string
	handler echo.HandlerFunc
}

func NewRoute(method string, path string, handler echo.HandlerFunc) Route {
	return &routeImpl{
		method:  method,
		path:    path,
		handler: handler,
	}
}

func (r *routeImpl) Method() string {
	return r.method
}

func (r *routeImpl) Register(path string, e *echo.Echo) {
	path = concatenateEndpoints(path, r.path)

	switch r.method {
	case "GET":
		registerGetRoute(path, r.handler, e)
	case "POST":
		registerPostRoute(path, r.handler, e)
	case "DELETE":
		registerDeleteRoute(path, r.handler, e)
	}
}
