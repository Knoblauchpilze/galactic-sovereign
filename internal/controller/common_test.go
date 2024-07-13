package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var errDefault = fmt.Errorf("some error")

type HandlerTestSuite[Service any] struct {
	suite.Suite

	generateTestFunc func(func(echo.Context, Service) error) echo.HandlerFunc
}

func (s *HandlerTestSuite[any]) TestCallsHandler() {
	assert := assert.New(s.T())

	handlerCalled := false
	in := func(_ echo.Context, _ any) error {
		handlerCalled = true
		return nil
	}

	h := s.generateTestFunc(in)

	err := h(dummyEchoContext())

	assert.Nil(err)
	assert.True(handlerCalled)
}

func (s *HandlerTestSuite[any]) TestPropagatesError() {
	assert := assert.New(s.T())

	in := func(_ echo.Context, _ any) error {
		return errDefault
	}

	h := s.generateTestFunc(in)

	err := h(dummyEchoContext())

	assert.Equal(errDefault, err)
}

func TestFromAuthServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.AuthService]{
		generateTestFunc: func(in func(echo.Context, service.AuthService) error) echo.HandlerFunc {
			return fromAuthServiceAwareHttpHandler(in, &mockAuthService{})
		},
	}

	suite.Run(t, &s)
}

func TestFromDbAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[db.ConnectionPool]{
		generateTestFunc: func(in func(echo.Context, db.ConnectionPool) error) echo.HandlerFunc {
			return fromDbAwareHttpHandler(in, &mockConnectionPool{})
		},
	}

	suite.Run(t, &s)
}

func TestFromPlanetServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.PlanetService]{
		generateTestFunc: func(in func(echo.Context, service.PlanetService) error) echo.HandlerFunc {
			return fromPlanetServiceAwareHttpHandler(in, &mockPlanetService{})
		},
	}

	suite.Run(t, &s)
}

func TestFromPlayerServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.PlayerService]{
		generateTestFunc: func(in func(echo.Context, service.PlayerService) error) echo.HandlerFunc {
			return fromPlayerServiceAwareHttpHandler(in, &mockPlayerService{})
		},
	}

	suite.Run(t, &s)
}

func TestFromUniverseServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.UniverseService]{
		generateTestFunc: func(in func(echo.Context, service.UniverseService) error) echo.HandlerFunc {
			return fromUniverseServiceAwareHttpHandler(in, &mockUniverseService{})
		},
	}

	suite.Run(t, &s)
}

func TestFromUserServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.UserService]{
		generateTestFunc: func(in func(echo.Context, service.UserService) error) echo.HandlerFunc {
			return fromUserServiceAwareHttpHandler(in, &mockUserService{})
		},
	}

	suite.Run(t, &s)
}

func dummyEchoContext() echo.Context {
	ctx, _ := generateTestEchoContextWithMethod(http.MethodGet)
	return ctx
}

func generateTestEchoContextWithMethod(method string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, "/", nil)
	return generateTestEchoContextFromRequest(req)
}

func generateTestEchoContextFromRequest(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)
	return ctx, rw
}
