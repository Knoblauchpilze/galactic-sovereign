package rest

import (
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockRoute struct {
	method string

	generatePathCalled int
	endpoint           string
}

type mockApiKeyRepository struct {
	repositories.ApiKeyRepository
}

func TestServer_Register_UsesPathFromRoute(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}

	NewServer(Config{}, &mockApiKeyRepository{}).Register(mr)
	assert.Equal(1, mr.generatePathCalled)
}

func TestServer_Register_PropagatesPathFromConfig(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}
	c := Config{
		Endpoint: "some-endpoint",
	}

	NewServer(c, &mockApiKeyRepository{}).Register(mr)
	assert.Equal(c.Endpoint, mr.endpoint)
}

func TestServer_Register_SanitizesPath(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}
	c := Config{
		Endpoint: "some-endpoint/",
	}

	NewServer(c, &mockApiKeyRepository{}).Register(mr)
	assert.Equal("some-endpoint", mr.endpoint)
}

func TestServer_Register_SupportsPost(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{
		method: http.MethodPost,
	}

	err := NewServer(Config{}, &mockApiKeyRepository{}).Register(mr)
	assert.Nil(err)
}

func TestServer_Register_SupportsGet(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{
		method: http.MethodGet,
	}

	err := NewServer(Config{}, &mockApiKeyRepository{}).Register(mr)
	assert.Nil(err)
}

func TestServer_Register_SupportsPatch(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{
		method: http.MethodPatch,
	}

	err := NewServer(Config{}, &mockApiKeyRepository{}).Register(mr)
	assert.Nil(err)
}

func TestServer_Register_SupportsDelete(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{
		method: http.MethodDelete,
	}

	err := NewServer(Config{}, &mockApiKeyRepository{}).Register(mr)
	assert.Nil(err)
}

func TestServer_Register_FailsForUnsupportedMethod(t *testing.T) {
	assert := assert.New(t)

	testMethods := []string{
		http.MethodPut,
		"not-a-http-method",
	}

	for _, method := range testMethods {
		t.Run(method, func(t *testing.T) {
			mr := &mockRoute{
				method: method,
			}

			err := NewServer(Config{}, &mockApiKeyRepository{}).Register(mr)
			assert.True(errors.IsErrorWithCode(err, UnsupportedMethod))
		})
	}

}

func (m *mockRoute) Method() string {
	return m.method
}

func (m *mockRoute) Handler() echo.HandlerFunc {
	return defaultHandler
}

func (m *mockRoute) GeneratePath(endpoint string) string {
	m.generatePathCalled++
	m.endpoint = endpoint
	return ""
}
