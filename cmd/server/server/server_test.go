package server

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockRoute struct {
	method string

	registerCalled int
	path           string
}

func TestServer_Register_DelegatesToRoute(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}

	s := New(Config{})
	s.Register(mr)

	assert.Equal(1, mr.registerCalled)
}

func TestServer_Register_SanitizesPath(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}
	c := Config{
		Endpoint: "some-endpoint/",
	}

	s := New(c)
	s.Register(mr)

	assert.Equal("some-endpoint", mr.path)
}

func TestServer_Register_UsesPathFromConfig(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}
	c := Config{
		Endpoint: "some-endpoint",
	}

	s := New(c)
	s.Register(mr)

	assert.Equal(c.Endpoint, mr.path)
}

func (m *mockRoute) Method() string {
	return m.method
}

func (m *mockRoute) Register(path string, e *echo.Echo) {
	m.registerCalled++
	m.path = path
}
