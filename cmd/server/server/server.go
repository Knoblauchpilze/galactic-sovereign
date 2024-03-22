package server

import (
	"fmt"
	"strings"

	"github.com/KnoblauchPilze/user-service/cmd/server/routes"
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type Server interface {
	Start() error
	Register(route routes.Route)
}

type serverImpl struct {
	endpoint   string
	port       uint16
	echoServer *echo.Echo
}

func New(conf Config) Server {
	return &serverImpl{
		endpoint:   strings.TrimSuffix(conf.Endpoint, "/"),
		port:       conf.Port,
		echoServer: createEchoContext(),
	}
}

func (s *serverImpl) Start() error {
	address := fmt.Sprintf(":%d", s.port)

	s.echoServer.Logger.Infof("Starting server at %s%s", s.endpoint, address)
	return s.echoServer.Start(address)
}

func (s *serverImpl) Register(route routes.Route) {
	path := concatenateEndpoints(s.endpoint, route.Path)

	if route.GetRoute != nil {
		getPath := concatenateEndpoints(path, ":id")
		s.echoServer.Logger.Debugf("Adding route GET %s", getPath)
		s.echoServer.GET(getPath, route.GetRoute)
	}
	if route.ListRoute != nil {
		s.echoServer.Logger.Debugf("Adding route GET %s", path)
		s.echoServer.GET(path, route.ListRoute)
	}
	if route.PostRoute != nil {
		s.echoServer.Logger.Debugf("Adding route POST %s", path)
		s.echoServer.POST(path, route.PostRoute)
	}
	if route.DeleteRoute != nil {
		s.echoServer.Logger.Debugf("Adding route DELETE %s", path)
		s.echoServer.DELETE(path, route.DeleteRoute)
	}
}

func createEchoContext() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = logger.New("server")

	e.Use(middleware.RequestTiming())
	e.Use(middleware.ResponseEnvelope())
	e.Use(middleware.ErrorMiddleware())
	e.Use(middleware.Recover())

	return e
}
