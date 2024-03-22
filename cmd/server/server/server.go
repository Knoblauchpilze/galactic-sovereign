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
	route.Register(s.endpoint, s.echoServer)
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
