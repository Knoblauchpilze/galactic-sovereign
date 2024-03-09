package routes

import (
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	mw "github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	e.Use(mw.RequestTiming())
	setupRecoverMiddleware(e)
	setupRoutes(e)

	return e
}

func setupRecoverMiddleware(e *echo.Echo) {
	recoverMiddleware := middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll:   true,
		DisablePrintStack: false,
	})

	e.Use(recoverMiddleware)
}
