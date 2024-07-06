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
	em "github.com/labstack/echo/v4/middleware"
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
		endpoint: strings.TrimSuffix(conf.BasePath, "/"),
		port:     conf.Port,

		server:       s,
		publicRoutes: publicRoutes,

		close: close,
	}
}

func NewServerWithApiKey(conf Config, apiKeyRepository repositories.ApiKeyRepository) Server {
	s := creationFunc()
	close := registerMiddlewares(s, conf.RateLimit)

	// https://github.com/labstack/echo/issues/1737#issuecomment-753355711
	publicRoutes := s.Group("")
	authorizedRoutes := s.Group("", middleware.ApiKey(apiKeyRepository))

	return &serverImpl{
		endpoint: strings.TrimSuffix(conf.BasePath, "/"),
		port:     conf.Port,

		server:           s,
		publicRoutes:     publicRoutes,
		authorizedRoutes: authorizedRoutes,

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
