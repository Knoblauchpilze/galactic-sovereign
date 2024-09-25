package game

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGameUpdateWatcher_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)

	called := false
	call := func(c echo.Context) error {
		called = true
		return nil
	}

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := GameUpdateWatcher(m, call)

	err := callable(ctx)

	assert.Nil(err)
	assert.True(called)
}

func TestGameUpdateWatcher_SchedulesActions(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := GameUpdateWatcher(m, defaultHandler)

	callable(ctx)

	assert.Equal(1, m.processActionsCalled)
}

func TestGameUpdateWatcher_ScheduleTimeIsAtTheMomentOfTheCall(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := GameUpdateWatcher(m, defaultHandler)

	beforeCall := time.Now()

	callable(ctx)

	assert.True(beforeCall.Before(m.until))
}

func TestGameUpdateWatcher_WhenServiceFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{
		err: errDefault,
	}
	ctx, _, rw := generateTestEchoContext()
	callable := GameUpdateWatcher(m, defaultHandler)

	err := callable(ctx)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, rw.Code)
	assert.Equal("\"Failed to process actions\"\n", rw.Body.String())
}

func TestGameUpdateWatcher_WhenServiceFails_DoesNotCallHandler(t *testing.T) {
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
	callable := GameUpdateWatcher(m, call)

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
