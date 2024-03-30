package rest

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var defaultHandler = func(c echo.Context) error { return nil }

func TestRoute_Method(t *testing.T) {
	assert := assert.New(t)

	r := NewRoute(http.MethodGet, "", defaultHandler)
	assert.Equal(http.MethodGet, r.Method())
}

func TestRoute_WithResource_Method(t *testing.T) {
	assert := assert.New(t)

	r := NewResourceRoute(http.MethodGet, "", defaultHandler)
	assert.Equal(http.MethodGet, r.Method())
}

type mockEchoContext struct {
	echo.Context
}

func TestRoute_Handler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	r := NewRoute(http.MethodGet, "", handler)
	actual := r.Handler()
	actual(mockEchoContext{})

	assert.True(handlerCalled)
}

func TestRoute_WithResource_Handler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return nil
	}

	r := NewResourceRoute(http.MethodGet, "", handler)
	actual := r.Handler()
	actual(mockEchoContext{})

	assert.True(handlerCalled)
}

type testCase struct {
	endpoint string
	path     string
	expected string
}

func TestRoute_GeneratePath(t *testing.T) {
	assert := assert.New(t)

	tests := []testCase{
		{endpoint: "endpoint", path: "path", expected: "/endpoint/path"},
		{endpoint: "endpoint", path: "/path", expected: "/endpoint/path"},
		{endpoint: "endpoint", path: "/path/", expected: "/endpoint/path"},
		{endpoint: "/endpoint", path: "path", expected: "/endpoint/path"},
		{endpoint: "/endpoint", path: "/path", expected: "/endpoint/path"},
		{endpoint: "/endpoint", path: "/path/", expected: "/endpoint/path"},
		{endpoint: "/endpoint/", path: "path", expected: "/endpoint/path"},
		{endpoint: "/endpoint/", path: "/path", expected: "/endpoint/path"},
		{endpoint: "/endpoint/", path: "/path/", expected: "/endpoint/path"},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {

			r := NewRoute(http.MethodGet, tc.path, defaultHandler)
			assert.Equal(tc.expected, r.GeneratePath(tc.endpoint))
		})
	}
}

func TestRoute_WithResource_GeneratePath(t *testing.T) {
	assert := assert.New(t)

	tests := []testCase{
		{endpoint: "endpoint", path: "path", expected: "/endpoint/path/:id"},
		{endpoint: "endpoint", path: "/path", expected: "/endpoint/path/:id"},
		{endpoint: "endpoint", path: "/path/", expected: "/endpoint/path/:id"},
		{endpoint: "/endpoint", path: "path", expected: "/endpoint/path/:id"},
		{endpoint: "/endpoint", path: "/path", expected: "/endpoint/path/:id"},
		{endpoint: "/endpoint", path: "/path/", expected: "/endpoint/path/:id"},
		{endpoint: "/endpoint/", path: "path", expected: "/endpoint/path/:id"},
		{endpoint: "/endpoint/", path: "/path", expected: "/endpoint/path/:id"},
		{endpoint: "/endpoint/", path: "/path/", expected: "/endpoint/path/:id"},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {

			r := NewResourceRoute(http.MethodGet, tc.path, defaultHandler)
			assert.Equal(tc.expected, r.GeneratePath(tc.endpoint))
		})
	}
}
