package rest

import "github.com/labstack/echo/v4"

type Route interface {
	Method() string
	Authorized() bool
	Handler() echo.HandlerFunc
	Path() string
}

type Routes []Route

type routeImpl struct {
	method      string
	authorized  bool
	path        string
	addIdInPath bool
	handler     echo.HandlerFunc
}

func NewRoute(method string, authorized bool, path string, handler echo.HandlerFunc) Route {
	return &routeImpl{
		method:      method,
		authorized:  authorized,
		path:        sanitizePath(path),
		addIdInPath: false,
		handler:     handler,
	}
}

func NewResourceRoute(method string, authorized bool, path string, handler echo.HandlerFunc) Route {
	return &routeImpl{
		method:      method,
		authorized:  authorized,
		path:        sanitizePath(path),
		addIdInPath: true,
		handler:     handler,
	}
}

func (r *routeImpl) Method() string {
	return r.method
}

func (r *routeImpl) Authorized() bool {
	return r.authorized
}

func (r *routeImpl) Handler() echo.HandlerFunc {
	return r.handler
}

func (r *routeImpl) Path() string {
	path := r.path
	if r.addIdInPath {
		path = concatenateEndpoints(path, ":id")
	}

	return path
}
