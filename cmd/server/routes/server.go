package routes

import (
	"os"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
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
	e.Logger = logger.New()

	setupLogMiddleware(e)
	setupRecoverMiddleware(e)
	setupRoutes(e)

	return e
}

func setupLogMiddleware(e *echo.Echo) {
	logMiddleware := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_rfc3339_nano} ${method} ${uri} -> ${status}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           os.Stdout,
	})

	e.Use(logMiddleware)
}

func setupRecoverMiddleware(e *echo.Echo) {
	recoverMiddleware := middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll:   true,
		DisablePrintStack: false,
	})

	e.Use(recoverMiddleware)
}
