package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
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

func TestUnit_FromBuildingActionServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.BuildingActionService]{
		generateTestFunc: func(in func(echo.Context, service.BuildingActionService) error) echo.HandlerFunc {
			return fromBuildingActionServiceAwareHttpHandler(in, &mockBuildingActionService{})
		},
	}

	suite.Run(t, &s)
}

func TestUnit_FromDbAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[db.ConnectionPool]{
		generateTestFunc: func(in func(echo.Context, db.ConnectionPool) error) echo.HandlerFunc {
			return fromDbAwareHttpHandler(in, &mockConnectionPool{})
		},
	}

	suite.Run(t, &s)
}

func TestUnit_FromPlanetServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.PlanetService]{
		generateTestFunc: func(in func(echo.Context, service.PlanetService) error) echo.HandlerFunc {
			return fromPlanetServiceAwareHttpHandler(in, &mockPlanetService{})
		},
	}

	suite.Run(t, &s)
}

func TestUnit_FromPlayerServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.PlayerService]{
		generateTestFunc: func(in func(echo.Context, service.PlayerService) error) echo.HandlerFunc {
			return fromPlayerServiceAwareHttpHandler(in, &mockPlayerService{})
		},
	}

	suite.Run(t, &s)
}

func TestUnit_FromUniverseServiceAwareHttpHandler(t *testing.T) {
	s := HandlerTestSuite[service.UniverseService]{
		generateTestFunc: func(in func(echo.Context, service.UniverseService) error) echo.HandlerFunc {
			return fromUniverseServiceAwareHttpHandler(in, &mockUniverseService{})
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

func generateTestEchoContextWithMethodAndId(method string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, "/", nil)
	ctx, rw := generateTestEchoContextFromRequest(req)

	ctx.SetParamNames("id")
	ctx.SetParamValues(defaultUuid.String())

	return ctx, rw
}

func generateTestEchoContextFromRequest(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)
	return ctx, rw
}
