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

type serverImpl struct {
	endpoint   string
	port       uint16
	echoServer *echo.Echo

	wg    sync.WaitGroup
	close chan bool
	err   error
}

func NewServer(conf Config, apiKeyRepository repositories.ApiKeyRepository) Server {
	ctx, close := createContextAndMiddlewares(apiKeyRepository)

	return &serverImpl{
		endpoint:   strings.TrimSuffix(conf.Endpoint, "/"),
		port:       conf.Port,
		echoServer: ctx,
		close:      close,
	}
}

func (s *serverImpl) Start() {
	address := fmt.Sprintf(":%d", s.port)

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		s.echoServer.Logger.Infof("Starting server at %s%s", s.endpoint, address)
		s.err = s.echoServer.Start(address)
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

func createContextAndMiddlewares(apiKeyRepository repositories.ApiKeyRepository) (*echo.Echo, chan bool) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = logger.New("server")

	e.Use(middleware.RequestTiming())
	e.Use(middleware.ResponseEnvelope())
	handler, close := middleware.ThrottleMiddleware(4, 2, 4)
	e.Use(handler)
	e.Use(middleware.ErrorMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.ApiKeyMiddleware(apiKeyRepository))

	return e, close
}
