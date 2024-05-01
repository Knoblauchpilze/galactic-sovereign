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
)

type Server interface {
	Start()
	Wait() error
	Stop()
	Register(route Route) error
}

type serverImpl struct {
	endpoint string
	port     uint16
	server   echoServer

	wg    sync.WaitGroup
	close chan bool
	err   error
}

var creationFunc = createEchoServer

func NewServer(conf Config, apiKeyRepository repositories.ApiKeyRepository) Server {
	s := creationFunc()
	close := registerMiddlewares(s, conf.RateLimit, apiKeyRepository)

	return &serverImpl{
		endpoint: strings.TrimSuffix(conf.BasePath, "/"),
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
	path := route.Path()
	path = concatenateEndpoints(s.endpoint, path)

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

func registerMiddlewares(server echoServer, rateLimit int, apiKeyRepository repositories.ApiKeyRepository) chan bool {
	server.Use(middleware.RequestTiming())
	server.Use(middleware.ResponseEnvelope())

	handler, close := middleware.ThrottleMiddleware(rateLimit, rateLimit, rateLimit)
	server.Use(handler)

	server.Use(middleware.ErrorMiddleware())
	server.Use(middleware.Recover())
	server.Use(middleware.ApiKeyMiddleware(apiKeyRepository))

	return close
}
