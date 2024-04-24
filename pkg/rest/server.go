package rest

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/labstack/echo/v4"
)

type Server interface {
	Start()
	Wait() error
	Stop()
	Register(route Route) error
}

type serverFramework interface {
	GET(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	POST(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	DELETE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	PATCH(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route

	Use(...echo.MiddlewareFunc)

	Start(string) error
}

type serverImpl struct {
	endpoint string
	port     uint16
	server   serverFramework

	wg    sync.WaitGroup
	close chan bool
	err   error
}

var creationFunc = createServerFramework

func NewServer(conf Config, apiKeyRepository repositories.ApiKeyRepository) Server {
	s := creationFunc()
	close := registerMiddlewares(s, apiKeyRepository)

	return &serverImpl{
		endpoint: strings.TrimSuffix(conf.Endpoint, "/"),
		port:     conf.Port,
		server:   s,
		close:    close,
	}
}

func (s *serverImpl) Start() {
	address := fmt.Sprintf(":%d", s.port)

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		logger.Infof("Starting server at %s%s", s.endpoint, address)
		s.err = s.server.Start(address)
	}()
}

func (s *serverImpl) Wait() error {
	s.wg.Wait()
	s.Stop()

	return s.err
}

func (s *serverImpl) Stop() {
	s.close <- true
}

func (s *serverImpl) Register(route Route) error {
	path := route.GeneratePath(s.endpoint)

	switch route.Method() {
	case http.MethodGet:
		s.server.GET(path, route.Handler())
	case http.MethodPost:
		s.server.POST(path, route.Handler())
	case http.MethodDelete:
		s.server.DELETE(path, route.Handler())
	case http.MethodPatch:
		s.server.PATCH(path, route.Handler())
	default:
		return errors.NewCode(UnsupportedMethod)
	}

	logger.Debugf("Registered %s %s", route.Method(), path)

	return nil
}

func createServerFramework() serverFramework {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = logger.New("server")
	return e
}

func registerMiddlewares(server serverFramework, apiKeyRepository repositories.ApiKeyRepository) chan bool {
	server.Use(middleware.RequestTiming())
	server.Use(middleware.ResponseEnvelope())
	handler, close := middleware.ThrottleMiddleware(4, 2, 4)
	server.Use(handler)
	server.Use(middleware.ErrorMiddleware())
	server.Use(middleware.Recover())
	server.Use(middleware.ApiKeyMiddleware(apiKeyRepository))

	return close
}
