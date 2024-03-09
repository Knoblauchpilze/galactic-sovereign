package routes

import (
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/labstack/echo"
)

type Server struct {
	echoServer *echo.Echo
}

func NewServer() *Server {

	s := &Server{
		echoServer: createEchoContext(),
	}

	return s
}

func (s *Server) Start(address string) {

	s.echoServer.Logger.Fatal(s.echoServer.Start(":1323"))
}

func createEchoContext() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = logger.New("server")

	e.Use(middleware.RequestTiming())
	e.Use(middleware.Recover())

	setupRoutes(e)

	return e
}
