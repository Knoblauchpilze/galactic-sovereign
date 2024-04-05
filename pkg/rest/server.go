package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/labstack/echo/v4"
)

type Server interface {
	Start() error
	Register(route Route) error
}

type serverImpl struct {
	endpoint   string
	port       uint16
	echoServer *echo.Echo
}

func NewServer(conf Config, apiKeyRepository repositories.ApiKeyRepository) Server {
	return &serverImpl{
		endpoint:   strings.TrimSuffix(conf.Endpoint, "/"),
		port:       conf.Port,
		echoServer: createEchoContext(apiKeyRepository),
	}
}

func (s *serverImpl) Start() error {
	address := fmt.Sprintf(":%d", s.port)

	s.echoServer.Logger.Infof("Starting server at %s%s", s.endpoint, address)
	return s.echoServer.Start(address)
}

func (s *serverImpl) Register(route Route) error {
	path := route.GeneratePath(s.endpoint)

	switch route.Method() {
	case http.MethodGet:
		s.echoServer.GET(path, route.Handler())
	case http.MethodPost:
		s.echoServer.POST(path, route.Handler())
	case http.MethodDelete:
		s.echoServer.DELETE(path, route.Handler())
	case http.MethodPatch:
		s.echoServer.PATCH(path, route.Handler())
	default:
		return errors.NewCode(UnsupportedMethod)
	}

	s.echoServer.Logger.Debugf("Registered %s %s", route.Method(), path)

	return nil
}

func createEchoContext(apiKeyRepository repositories.ApiKeyRepository) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = logger.New("server")

	e.Use(middleware.RequestTiming())
	e.Use(middleware.ResponseEnvelope())
	e.Use(middleware.ErrorMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.ApiKeyMiddleware(apiKeyRepository))

	return e
}
