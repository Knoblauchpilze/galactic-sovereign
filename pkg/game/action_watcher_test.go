package game

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

var defaultHandler = func(c echo.Context) error {
	return nil
}

type mockActionService struct {
	ActionService

	err error

	processActionsCalled int
	until                time.Time
}

func (m *mockActionService) ProcessActionsUntil(ctx context.Context, until time.Time) error {
	m.processActionsCalled++
	m.until = until
	return m.err
}

func TestActionWatcher_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)

	called := false
	call := func(c echo.Context) error {
		called = true
		return nil
	}

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := ActionWatcher(m, call)

	err := callable(ctx)

	assert.Nil(err)
	assert.True(called)
}

func TestActionWatcher_SchedulesActions(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := ActionWatcher(m, defaultHandler)

	callable(ctx)

	assert.Equal(1, m.processActionsCalled)
}

func TestActionWatcher_ScheduleTimeIsAtTheMomentOfTheCall(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{}
	ctx, _, _ := generateTestEchoContext()
	callable := ActionWatcher(m, defaultHandler)

	beforeCall := time.Now()

	callable(ctx)

	assert.True(beforeCall.Before(m.until))
}

func TestActionWatcher_WhenServiceFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{
		err: errDefault,
	}
	ctx, _, rw := generateTestEchoContext()
	callable := ActionWatcher(m, defaultHandler)

	err := callable(ctx)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, rw.Code)
	assert.Equal("\"Failed to process action\"\n", rw.Body.String())
}

func TestActionWatcher_WhenServiceFails_DoesNotCallHandler(t *testing.T) {
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
	callable := ActionWatcher(m, call)

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
