package rest

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockRoute struct {
	method string

	generatePathCalled int
	basePath           string
}

type mockApiKeyRepository struct {
	repositories.ApiKeyRepository
}

type mockServerFramework struct {
	getCalled    int
	postCalled   int
	deleteCalled int
	patchCalled  int

	useCalled int

	startCalled int

	address     string
	path        string
	middlewares []echo.MiddlewareFunc

	sleep time.Duration
	err   error
}

var errDefault = fmt.Errorf("some error")

func TestServer_Register_UsesPathFromRoute(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}

	s := NewServer(Config{}, &mockApiKeyRepository{})
	defer s.Stop()

	s.Register(mr)
	assert.Equal(1, mr.generatePathCalled)
}

func TestServer_Register_PropagatesPathFromConfig(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}
	c := Config{
		BasePath: "some-path",
	}

	s := NewServer(c, &mockApiKeyRepository{})
	defer s.Stop()

	s.Register(mr)
	assert.Equal(c.BasePath, mr.basePath)
}

func TestServer_Register_SanitizesPath(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}
	c := Config{
		BasePath: "some-path/",
	}

	s := NewServer(c, &mockApiKeyRepository{})
	defer s.Stop()

	s.Register(mr)
	assert.Equal("some-path", mr.basePath)
}

func TestServer_Register_SupportsPost(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	mr := &mockRoute{
		method: http.MethodPost,
	}

	s := NewServer(Config{}, &mockApiKeyRepository{})
	defer s.Stop()

	err := s.Register(mr)
	assert.Nil(err)
	assert.Equal(1, ms.postCalled)
}

func TestServer_Register_SupportsGet(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	mr := &mockRoute{
		method: http.MethodGet,
	}

	s := NewServer(Config{}, &mockApiKeyRepository{})
	defer s.Stop()

	err := s.Register(mr)
	assert.Nil(err)
	assert.Equal(1, ms.getCalled)
}

func TestServer_Register_SupportsPatch(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	mr := &mockRoute{
		method: http.MethodPatch,
	}

	s := NewServer(Config{}, &mockApiKeyRepository{})
	defer s.Stop()

	err := s.Register(mr)
	assert.Nil(err)
	assert.Equal(1, ms.patchCalled)
}

func TestServer_Register_SupportsDelete(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	mr := &mockRoute{
		method: http.MethodDelete,
	}

	s := NewServer(Config{}, &mockApiKeyRepository{})
	defer s.Stop()

	err := s.Register(mr)
	assert.Nil(err)
	assert.Equal(1, ms.deleteCalled)
}

func TestServer_Register_FailsForUnsupportedMethod(t *testing.T) {
	assert := assert.New(t)

	testMethods := []string{
		http.MethodPut,
		"not-a-http-method",
	}

	for _, method := range testMethods {
		t.Run(method, func(t *testing.T) {
			t.Cleanup(resetCreatorFunc)

			setupMockServer()

			mr := &mockRoute{
				method: method,
			}

			s := NewServer(Config{}, &mockApiKeyRepository{})
			defer s.Stop()

			err := s.Register(mr)
			assert.True(errors.IsErrorWithCode(err, UnsupportedMethod))
		})
	}
}

func TestServer_Start_CallsServerStart(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()

	s := NewServer(Config{}, mockApiKeyRepository{})
	s.Start()
	s.Wait()

	assert.Equal(1, ms.startCalled)
}

func TestServer_Start_UsesPortFromConfig(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	conf := Config{
		Port: 36,
	}

	s := NewServer(conf, mockApiKeyRepository{})
	s.Start()
	s.Wait()

	assert.Equal(":36", ms.address)
}

func TestServer_Start_ReturnsServerError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	ms.err = errDefault

	s := NewServer(Config{}, mockApiKeyRepository{})
	s.Start()
	err := s.Wait()

	assert.Equal(errDefault, err)
}

func TestServer_Wait_WhenNotStarted_ReturnsImmediately(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	setupMockServer()

	s := NewServer(Config{}, mockApiKeyRepository{})
	err := s.Wait()

	assert.Nil(err)
}

func TestServer_Wait_WhenStarted_TakesTime(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()

	ms.sleep = time.Second

	s := NewServer(Config{}, mockApiKeyRepository{})
	start := time.Now()
	s.Wait()
	elapsed := time.Since(start)

	assert.GreaterOrEqual(900*time.Millisecond, elapsed)
}

func TestRegisterMiddlewares_registersExpectedMiddlewareCount(t *testing.T) {
	assert := assert.New(t)

	ms := mockServerFramework{}

	c := registerMiddlewares(&ms, 1, mockApiKeyRepository{})
	defer func() {
		c <- true
	}()

	assert.Equal(6, len(ms.middlewares))
}

func setupMockServer() *mockServerFramework {
	server := &mockServerFramework{}

	creationFunc = func() echoServer {
		return server
	}

	return server
}

func resetCreatorFunc() {
	creationFunc = createEchoServer
}

func (m *mockRoute) Method() string {
	return m.method
}

func (m *mockRoute) Authorized() bool {
	return false
}

func (m *mockRoute) Handler() echo.HandlerFunc {
	return defaultHandler
}

func (m *mockRoute) GeneratePath(basePath string) string {
	m.generatePathCalled++
	m.basePath = basePath
	return ""
}

func (m *mockServerFramework) GET(path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *echo.Route {
	m.getCalled++
	m.path = path
	return nil
}

func (m *mockServerFramework) POST(path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *echo.Route {
	m.postCalled++
	m.path = path
	return nil
}

func (m *mockServerFramework) DELETE(path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *echo.Route {
	m.deleteCalled++
	m.path = path
	return nil
}

func (m *mockServerFramework) PATCH(path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *echo.Route {
	m.patchCalled++
	m.path = path
	return nil
}

func (m *mockServerFramework) Group(prefix string, middlewares ...echo.MiddlewareFunc) *echo.Group {
	return nil
}

func (m *mockServerFramework) Use(middlewares ...echo.MiddlewareFunc) {
	m.useCalled++
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *mockServerFramework) Start(address string) error {
	m.startCalled++
	m.address = address
	time.Sleep(m.sleep)
	return m.err
}
