package game

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")
var defaultHandler = func(c echo.Context) error { return nil }

type mockActionService struct {
	ActionService

	err error

	processActionsCalled int
	planet               uuid.UUID
	until                time.Time
}

type mockPlanetResourceService struct {
	PlanetResourceService

	err error

	updatePlanetUntilCalled int
	planet                  uuid.UUID
	until                   time.Time
}

func TestUnit_Route_Method(t *testing.T) {
	assert := assert.New(t)

	r := NewRoute(http.MethodGet, "", defaultHandler, &mockActionService{}, &mockPlanetResourceService{})
	assert.Equal(http.MethodGet, r.Method())
}

func TestUnit_Route_WithResource_Method(t *testing.T) {
	assert := assert.New(t)

	r := NewResourceRoute(http.MethodGet, "", defaultHandler, &mockActionService{}, &mockPlanetResourceService{})
	assert.Equal(http.MethodGet, r.Method())
}

func TestUnit_Route_Authorized(t *testing.T) {
	assert := assert.New(t)

	public := NewRoute(http.MethodGet, "", defaultHandler, &mockActionService{}, &mockPlanetResourceService{})
	assert.Equal(false, public.Authorized())
}

func TestUnit_Route_WithResource_Authorized(t *testing.T) {
	assert := assert.New(t)

	public := NewResourceRoute(http.MethodGet, "", defaultHandler, &mockActionService{}, &mockPlanetResourceService{})
	assert.Equal(false, public.Authorized())
}

func TestUnit_Route_Handler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	r := NewRoute(http.MethodGet, "", handler, &mockActionService{}, &mockPlanetResourceService{})
	actual := r.Handler()
	err := actual(dummyEchoContext())

	assert.True(handlerCalled)
	assert.Nil(err)
}

func TestUnit_Route_WithResource_Handler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	r := NewResourceRoute(http.MethodGet, "", handler, &mockActionService{}, &mockPlanetResourceService{})
	actual := r.Handler()
	err := actual(dummyEchoContext())

	assert.True(handlerCalled)
	assert.Nil(err)
}

func TestUnit_Route_WhenNoPlanetId_DoesNotScheduleActions(t *testing.T) {
	assert := assert.New(t)

	handler := func(c echo.Context) error {
		return nil
	}

	m := &mockActionService{}

	r := NewResourceRoute(http.MethodGet, "", handler, m, &mockPlanetResourceService{})
	actual := r.Handler()
	ctx, _, _ := generateTestEchoContext()
	err := actual(ctx)

	assert.Nil(err)
	assert.Equal(0, m.processActionsCalled)
}

func TestUnit_Route_CallsActionService(t *testing.T) {
	assert := assert.New(t)

	handler := func(c echo.Context) error {
		return nil
	}

	m := &mockActionService{}

	r := NewResourceRoute(http.MethodGet, "", handler, m, &mockPlanetResourceService{})
	actual := r.Handler()
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	err := actual(ctx)

	assert.Nil(err)
	assert.Equal(1, m.processActionsCalled)
	assert.Equal(someUuid, m.planet)
}

func TestUnit_Route_ScheduleActionsIsAtTheRightTime(t *testing.T) {
	assert := assert.New(t)

	handler := func(c echo.Context) error {
		return nil
	}

	m := &mockActionService{}

	beforeCall := time.Now()

	r := NewResourceRoute(http.MethodGet, "", handler, m, &mockPlanetResourceService{})
	actual := r.Handler()
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	err := actual(ctx)

	assert.Nil(err)
	assert.True(beforeCall.Before(m.until))
}

func TestUnit_Route_WhenActionServiceFails_DoesNotCallHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	m := &mockActionService{
		err: errDefault,
	}

	r := NewResourceRoute(http.MethodGet, "", handler, m, &mockPlanetResourceService{})
	actual := r.Handler()
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	err := actual(ctx)

	assert.False(handlerCalled)
	assert.Nil(err)
}

func TestUnit_Route_WhenActionServiceFails_SetsStatusToInternalError(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{
		err: errDefault,
	}
	ctx, _, rw := generateTestEchoContextWithPlanetId()

	r := NewResourceRoute(http.MethodGet, "", defaultHandler, m, &mockPlanetResourceService{})
	actual := r.Handler()
	err := actual(ctx)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, rw.Code)
	assert.Equal("\"Failed to process actions\"\n", rw.Body.String())
}

func TestUnit_Route_WhenNoPlanetId_DoesNotUpdatePlanetResourceService(t *testing.T) {
	assert := assert.New(t)

	handler := func(c echo.Context) error {
		return nil
	}

	m := &mockPlanetResourceService{}

	r := NewResourceRoute(http.MethodGet, "", handler, &mockActionService{}, m)
	actual := r.Handler()
	ctx, _, _ := generateTestEchoContext()
	err := actual(ctx)

	assert.Nil(err)
	assert.Equal(0, m.updatePlanetUntilCalled)
}

func TestUnit_Route_CallsPlanetResourceService(t *testing.T) {
	assert := assert.New(t)

	handler := func(c echo.Context) error {
		return nil
	}

	m := &mockPlanetResourceService{}

	r := NewResourceRoute(http.MethodGet, "", handler, &mockActionService{}, m)
	actual := r.Handler()
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	err := actual(ctx)

	assert.Nil(err)
	assert.Equal(1, m.updatePlanetUntilCalled)
	assert.Equal(someUuid, m.planet)
}

func TestUnit_Route_UpdateResourcesIsAtTheRightTime(t *testing.T) {
	assert := assert.New(t)

	handler := func(c echo.Context) error {
		return nil
	}

	m := &mockPlanetResourceService{}

	beforeCall := time.Now()

	r := NewResourceRoute(http.MethodGet, "", handler, &mockActionService{}, m)
	actual := r.Handler()
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	err := actual(ctx)

	assert.Nil(err)
	assert.True(beforeCall.Before(m.until))
}

func TestUnit_Route_WhenPlanetResourceServiceFails_DoesNotCallHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	m := &mockPlanetResourceService{
		err: errDefault,
	}

	r := NewResourceRoute(http.MethodGet, "", handler, &mockActionService{}, m)
	actual := r.Handler()
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	err := actual(ctx)

	assert.False(handlerCalled)
	assert.Nil(err)
}

func TestUnit_Route_WhenPlanetResourceServiceFails_SetsStatusToInternalError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{
		err: errDefault,
	}
	ctx, _, rw := generateTestEchoContextWithPlanetId()

	r := NewResourceRoute(http.MethodGet, "", defaultHandler, &mockActionService{}, m)
	actual := r.Handler()
	err := actual(ctx)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, rw.Code)
	assert.Equal("\"Failed to update resources\"\n", rw.Body.String())
}

type routeTestCase struct {
	path     string
	expected string
}

func TestUnit_Route_Path(t *testing.T) {
	assert := assert.New(t)

	tests := []routeTestCase{
		{path: "path", expected: "/path"},
		{path: "/path", expected: "/path"},
		{path: "/path/", expected: "/path"},
		{path: "path/", expected: "/path"},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {

			r := NewRoute(http.MethodGet, tc.path, defaultHandler, &mockActionService{}, &mockPlanetResourceService{})
			assert.Equal(tc.expected, r.Path())
		})
	}
}

func TestUnit_Route_WithResource_GeneratePath(t *testing.T) {
	assert := assert.New(t)

	tests := []routeTestCase{
		{path: "path", expected: "/path/:id"},
		{path: "/path", expected: "/path/:id"},
		{path: "/path/", expected: "/path/:id"},
		{path: "path/", expected: "/path/:id"},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {

			r := NewResourceRoute(http.MethodGet, tc.path, defaultHandler, &mockActionService{}, &mockPlanetResourceService{})
			assert.Equal(tc.expected, r.Path())
		})
	}
}

func TestUnit_Route_WithResource_WhenIdPlaceHolderAlreadyExists_DoNotGeneratePath(t *testing.T) {
	assert := assert.New(t)

	path := "/path/:id/addendum"

	r := NewResourceRoute(http.MethodGet, path, defaultHandler, &mockActionService{}, &mockPlanetResourceService{})

	assert.Equal(path, r.Path())
}

func dummyEchoContext() echo.Context {
	e, _, _ := generateTestEchoContext()
	return e
}

func (m *mockActionService) ProcessActionsUntil(ctx context.Context, planet uuid.UUID, until time.Time) error {
	m.processActionsCalled++
	m.planet = planet
	m.until = until
	return m.err
}

func (m *mockPlanetResourceService) UpdatePlanetUntil(ctx context.Context, planet uuid.UUID, until time.Time) error {
	m.updatePlanetUntilCalled++
	m.planet = planet
	m.until = until
	return m.err
}
