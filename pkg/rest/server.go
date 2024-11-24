package rest

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/logger"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/middleware"
	em "github.com/labstack/echo/v4/middleware"
)

type Server interface {
	Start()
	Wait() error
	Stop() error
	Register(route Route) error
}

type serverImpl struct {
	endpoint string
	port     uint16

	server           echoServer
	publicRoutes     echoRouter
	authorizedRoutes echoRouter

	wg    sync.WaitGroup
	close chan bool
	err   error
}

var creationFunc = createEchoServerWrapper

func NewServer(conf Config) Server {
	s := creationFunc()
	close := registerMiddlewares(s, conf.RateLimit)

	publicRoutes := s.Group("")

	return &serverImpl{
		endpoint: ConcatenateEndpoints(conf.BasePath, conf.Prefix),
		port:     conf.Port,

		server:       s,
		publicRoutes: publicRoutes,

		close: close,
	}
}

func (s *serverImpl) Start() {
	address := fmt.Sprintf(":%d", s.port)

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		logger.Infof("Starting server at %s for route %s", address, s.endpoint)
		s.err = s.server.Start(address)
		logger.Infof("Server at %s gracefully shutdown", address)
	}()
}

func (s *serverImpl) Wait() error {
	s.wg.Wait()
	s.Stop()

	return s.err
}

func (s *serverImpl) Stop() error {
	s.close <- true

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *serverImpl) Register(route Route) error {
	path := route.Path()
	path = ConcatenateEndpoints(s.endpoint, path)

	router := s.publicRoutes
	if route.Authorized() {
		if s.authorizedRoutes == nil {
			return errors.NewCode(AuthorizationNotSupported)
		}
		router = s.authorizedRoutes
	}

	switch route.Method() {
	case http.MethodGet:
		router.GET(path, route.Handler())
	case http.MethodPost:
		router.POST(path, route.Handler())
	case http.MethodDelete:
		router.DELETE(path, route.Handler())
	case http.MethodPatch:
		router.PATCH(path, route.Handler())
	default:
		return errors.NewCode(UnsupportedMethod)
	}

	logger.Debugf("Registered %s %s", route.Method(), path)

	return nil
}

func registerMiddlewares(server echoServer, rateLimit int) chan bool {
	// https://stackoverflow.com/questions/74020538/cors-preflight-did-not-succeed
	// https://stackoverflow.com/questions/6660019/restful-api-methods-head-options
	corsConf := em.CORSConfig{
		// https://www.stackhawk.com/blog/golang-cors-guide-what-it-is-and-how-to-enable-it/
		// Same as the default value
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
		},
	}
	server.Use(em.CORSWithConfig(corsConf))
	server.Use(em.Gzip())

	server.Use(middleware.RequestTiming())
	server.Use(middleware.ResponseEnvelope())

	handler, close := middleware.Throttle(rateLimit, rateLimit, rateLimit)
	server.Use(handler)

	server.Use(middleware.Error())
	server.Use(middleware.Recover())

	return close
}
