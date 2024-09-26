package game

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var someUuid = uuid.MustParse("83a57c6c-7a5a-4c6a-87c9-2d6445f805a2")

func TestGameUpdateWatcher_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)

	called := false
	call := func(c echo.Context) error {
		called = true
		return nil
	}

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := GameUpdateWatcher(m, &mockPlanetResourceService{}, call)

	err := callable(ctx)

	assert.Nil(err)
	assert.True(called)
}

func TestGameUpdateWatcher_SchedulesActions(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := GameUpdateWatcher(m, &mockPlanetResourceService{}, defaultHandler)

	callable(ctx)

	assert.Equal(1, m.processActionsCalled)
}

func TestGameUpdateWatcher_ScheduleActionsTimeIsAtTheMomentOfTheCall(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := GameUpdateWatcher(m, &mockPlanetResourceService{}, defaultHandler)

	beforeCall := time.Now()

	callable(ctx)

	assert.True(beforeCall.Before(m.until))
}

func TestGameUpdateWatcher_WhenActionServiceFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{
		err: errDefault,
	}
	ctx, _, rw := generateTestEchoContext()
	callable := GameUpdateWatcher(m, &mockPlanetResourceService{}, defaultHandler)

	err := callable(ctx)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, rw.Code)
	assert.Equal("\"Failed to process actions\"\n", rw.Body.String())
}

func TestGameUpdateWatcher_WhenActionServiceFails_DoesNotCallHandler(t *testing.T) {
	assert := assert.New(t)

	called := false
	call := func(c echo.Context) error {
		called = true
		return nil
	}

	m := &mockActionService{
		err: errDefault,
	}
	ctx, _, _ := generateTestEchoContext()
	callable := GameUpdateWatcher(m, &mockPlanetResourceService{}, call)

	err := callable(ctx)

	assert.Nil(err)
	assert.False(called)
}

func TestGameUpdateWatcher_WhenNoPlanetId_DoesNotCallUpdateOfPlanetResource(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{}
	ctx, _, _ := generateTestEchoContext()
	callable := GameUpdateWatcher(&mockActionService{}, m, defaultHandler)

	callable(ctx)

	assert.Equal(0, m.updatePlanetUntilCalled)
}

func TestGameUpdateWatcher_SchedulesUpdateOfResources(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{}
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	callable := GameUpdateWatcher(&mockActionService{}, m, defaultHandler)

	callable(ctx)

	assert.Equal(1, m.updatePlanetUntilCalled)
	assert.Equal(someUuid, m.planet)
}

func TestGameUpdateWatcher_UpdateResourcesTimeIsAtTheMomentOfTheCall(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{}
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	callable := GameUpdateWatcher(&mockActionService{}, m, defaultHandler)

	beforeCall := time.Now()

	callable(ctx)

	assert.True(beforeCall.Before(m.until))
}

func TestGameUpdateWatcher_WhenPlanetResourceServiceFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{
		err: errDefault,
	}
	ctx, _, rw := generateTestEchoContextWithPlanetId()
	callable := GameUpdateWatcher(&mockActionService{}, m, defaultHandler)

	err := callable(ctx)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, rw.Code)
	assert.Equal("\"Failed to update resources\"\n", rw.Body.String())
}

func TestGameUpdateWatcher_WhenPlanetResourceServiceFails_DoesNotCallHandler(t *testing.T) {
	assert := assert.New(t)

	called := false
	call := func(c echo.Context) error {
		called = true
		return nil
	}

	m := &mockPlanetResourceService{
		err: errDefault,
	}
	ctx, _, _ := generateTestEchoContextWithPlanetId()
	callable := GameUpdateWatcher(&mockActionService{}, m, call)

	err := callable(ctx)

	assert.Nil(err)
	assert.False(called)
}

func generateTestEchoContext() (echo.Context, *http.Request, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)

	return ctx, req, rw
}

func generateTestEchoContextWithPlanetId() (echo.Context, *http.Request, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)

	ctx.SetParamNames("id")
	ctx.SetParamValues(someUuid.String())

	return ctx, req, rw
}
