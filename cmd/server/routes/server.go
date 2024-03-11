package routes

import (
	"fmt"
	"strings"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type Server struct {
	endpoint   string
	port       uint16
	echoServer *echo.Echo
}

func NewServer(endpoint string, port uint16) *Server {
	return &Server{
		endpoint:   strings.TrimSuffix(endpoint, "/"),
		port:       port,
		echoServer: createEchoContext(),
	}
}

func (s *Server) Start() {
	address := fmt.Sprintf(":%d", s.port)
	s.echoServer.Logger.Fatal(s.echoServer.Start(address))
}

func (s *Server) Register(route Route) {
	path := concatenateEndpoints(s.endpoint, route.Path)

	if route.GetRoute != nil {
		s.echoServer.GET(path, route.GetRoute)
	}
	if route.PostRoute != nil {
		s.echoServer.POST(path, route.PostRoute)
	}
	if route.DeleteRoute != nil {
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
