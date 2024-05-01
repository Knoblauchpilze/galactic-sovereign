package rest

import "github.com/labstack/echo/v4"

type Route interface {
	Method() string
	Handler() echo.HandlerFunc
	GeneratePath(basePath string) string
}

type Routes []Route

type routeImpl struct {
	method      string
	path        string
	addIdInPath bool
	handler     echo.HandlerFunc
}

func NewRoute(method string, path string, handler echo.HandlerFunc) Route {
	return &routeImpl{
		method:      method,
		path:        path,
		addIdInPath: false,
		handler:     handler,
	}
}

func NewResourceRoute(method string, path string, handler echo.HandlerFunc) Route {
	return &routeImpl{
		method:      method,
		path:        path,
		addIdInPath: true,
		handler:     handler,
	}
}

func (r *routeImpl) Method() string {
	return r.method
}

func (r *routeImpl) Handler() echo.HandlerFunc {
	return r.handler
}

func (r *routeImpl) GeneratePath(basePath string) string {
	path := concatenateEndpoints(basePath, r.path)
	if r.addIdInPath {
		path = concatenateEndpoints(path, ":id")
	}

	return path
}
