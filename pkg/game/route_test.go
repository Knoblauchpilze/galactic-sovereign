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
var defaultHandler = func(c echo.Context) error { return nil }

type mockActionService struct {
	ActionService

	err error

	processActionsCalled int
	until                time.Time
}

func TestRoute_Method(t *testing.T) {
	assert := assert.New(t)

	r := NewRoute(http.MethodGet, "", defaultHandler, &mockActionService{})
	assert.Equal(http.MethodGet, r.Method())
}

func TestRoute_WithResource_Method(t *testing.T) {
	assert := assert.New(t)

	r := NewResourceRoute(http.MethodGet, "", defaultHandler, &mockActionService{})
	assert.Equal(http.MethodGet, r.Method())
}

func TestRoute_Authorized(t *testing.T) {
	assert := assert.New(t)

	public := NewRoute(http.MethodGet, "", defaultHandler, &mockActionService{})
	assert.Equal(false, public.Authorized())
}

func TestRoute_WithResource_Authorized(t *testing.T) {
	assert := assert.New(t)

	public := NewResourceRoute(http.MethodGet, "", defaultHandler, &mockActionService{})
	assert.Equal(false, public.Authorized())
}

func TestRoute_Handler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	r := NewRoute(http.MethodGet, "", handler, &mockActionService{})
	actual := r.Handler()
	err := actual(dummyEchoContext())

	assert.True(handlerCalled)
	assert.Nil(err)
}

func TestRoute_WithResource_Handler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	r := NewResourceRoute(http.MethodGet, "", handler, &mockActionService{})
	actual := r.Handler()
	err := actual(dummyEchoContext())

	assert.True(handlerCalled)
	assert.Nil(err)
}

func TestRoute_WhenServiceFails_DoesNotCallHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	m := &mockActionService{
		err: errDefault,
	}

	r := NewResourceRoute(http.MethodGet, "", handler, m)
	actual := r.Handler()
	err := actual(dummyEchoContext())

	assert.False(handlerCalled)
	assert.Nil(err)
}

func TestRoute_WhenServiceFails_SetsStatusToInternalError(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{
		err: errDefault,
	}
	ctx, rw := dummyEchoContextWithRecorder()

	r := NewResourceRoute(http.MethodGet, "", defaultHandler, m)
	actual := r.Handler()
	err := actual(ctx)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, rw.Code)
	assert.Equal("\"Failed to process action\"\n", rw.Body.String())
}

type testCase struct {
	path     string
	expected string
}

func TestRoute_Path(t *testing.T) {
	assert := assert.New(t)

	tests := []testCase{
		{path: "path", expected: "/path"},
		{path: "/path", expected: "/path"},
		{path: "/path/", expected: "/path"},
		{path: "path/", expected: "/path"},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {

			r := NewRoute(http.MethodGet, tc.path, defaultHandler, &mockActionService{})
			assert.Equal(tc.expected, r.Path())
		})
	}
}

func TestRoute_WithResource_GeneratePath(t *testing.T) {
	assert := assert.New(t)

	tests := []testCase{
		{path: "path", expected: "/path/:id"},
		{path: "/path", expected: "/path/:id"},
		{path: "/path/", expected: "/path/:id"},
		{path: "path/", expected: "/path/:id"},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {

			r := NewResourceRoute(http.MethodGet, tc.path, defaultHandler, &mockActionService{})
			assert.Equal(tc.expected, r.Path())
		})
	}
}

func TestRoute_WithResource_WhenIdPlaceHolderAlreadyExists_DoNotGeneratePath(t *testing.T) {
	assert := assert.New(t)

	path := "/path/:id/addendum"

	r := NewResourceRoute(http.MethodGet, path, defaultHandler, &mockActionService{})

	assert.Equal(path, r.Path())
}

func dummyEchoContext() echo.Context {
	e, _ := dummyEchoContextWithRecorder()
	return e
}

func dummyEchoContextWithRecorder() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()

	return e.NewContext(req, rw), rw
}

func (m *mockActionService) ProcessActionsUntil(ctx context.Context, until time.Time) error {
	m.processActionsCalled++
	m.until = until
	return m.err
}
