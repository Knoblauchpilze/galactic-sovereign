package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockRoute struct {
	method     string
	path       string
	authorized bool

	pathCalled int
}

type groupData struct {
	prefix      string
	group       *mockEchoRouter
	middlewares int
}

type mockEchoServer struct {
	mockEchoRouter

	groupCalled int
	groups      []groupData

	startCalled    int
	shutdownCalled int

	sleep time.Duration
	err   error
}

type mockEchoRouter struct {
	getCalled    int
	postCalled   int
	deleteCalled int
	patchCalled  int

	useCalled int

	address     string
	path        string
	middlewares []echo.MiddlewareFunc
}

var errDefault = fmt.Errorf("some error")

func TestServer_DefinesOneGroup_WithoutMiddleware(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()

	s := NewServer(Config{})
	defer s.Stop()

	assert.Equal(1, ms.groupCalled)
	assert.Equal(0, ms.groups[0].middlewares)
}

func TestServer_Register_UsesPathFromRoute(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{}

	s := NewServer(Config{})
	defer s.Stop()

	s.Register(mr)
	assert.Equal(1, mr.pathCalled)
}

func TestServer_RegisterAuthorized_ExpectsErrorDueToMissingAuthorization(t *testing.T) {
	assert := assert.New(t)

	mr := &mockRoute{
		authorized: true,
	}

	s := NewServer(Config{})
	defer s.Stop()

	err := s.Register(mr)
	assert.True(errors.IsErrorWithCode(err, AuthorizationNotSupported))
}

func TestServer_Register_PropagatesPathFromConfig(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	mr := &mockRoute{
		method: http.MethodGet,
		path:   "path",
	}
	c := Config{
		BasePath: "some-path",
	}
	ms := setupMockServer()

	s := NewServer(c)
	defer s.Stop()

	s.Register(mr)
	assert.Equal("/some-path/path", ms.groups[0].group.path)
}

func TestServer_Register_SanitizesPath(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	mr := &mockRoute{
		method: http.MethodGet,
		path:   "/addition/",
	}
	c := Config{
		BasePath: "some-path/",
	}
	ms := setupMockServer()

	s := NewServer(c)
	defer s.Stop()

	s.Register(mr)
	assert.Equal("/some-path/addition", ms.groups[0].group.path)
}

func TestServer_Register_SupportsPost(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	mr := &mockRoute{
		method: http.MethodPost,
	}

	s := NewServer(Config{})
	defer s.Stop()

	err := s.Register(mr)
	assert.Nil(err)
	assert.Equal(1, ms.groups[0].group.postCalled)
}

func TestServer_Register_SupportsGet(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	mr := &mockRoute{
		method: http.MethodGet,
	}

	s := NewServer(Config{})
	defer s.Stop()

	err := s.Register(mr)
	assert.Nil(err)
	assert.Equal(1, ms.groups[0].group.getCalled)
}

func TestServer_Register_SupportsPatch(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	mr := &mockRoute{
		method: http.MethodPatch,
	}

	s := NewServer(Config{})
	defer s.Stop()

	err := s.Register(mr)
	assert.Nil(err)
	assert.Equal(1, ms.groups[0].group.patchCalled)
}

func TestServer_Register_SupportsDelete(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	mr := &mockRoute{
		method: http.MethodDelete,
	}

	s := NewServer(Config{})
	defer s.Stop()

	err := s.Register(mr)
	assert.Nil(err)
	assert.Equal(1, ms.groups[0].group.deleteCalled)
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

			s := NewServer(Config{})
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

	s := NewServer(Config{})
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

	s := NewServer(conf)
	s.Start()
	s.Wait()

	assert.Equal(":36", ms.address)
}

func TestServer_Start_ReturnsServerError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()
	ms.err = errDefault

	s := NewServer(Config{})
	s.Start()
	err := s.Wait()

	assert.Equal(errDefault, err)
}

func TestServer_Wait_WhenNotStarted_ReturnsImmediately(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	setupMockServer()

	s := NewServer(Config{})
	start := time.Now()
	err := s.Wait()
	elapsed := time.Since(start)

	assert.Nil(err)
	assert.True(elapsed < time.Duration(100*time.Millisecond))
}

func TestServer_Wait_WhenStarted_TakesTime(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetCreatorFunc)

	ms := setupMockServer()

	ms.sleep = time.Second

	s := NewServer(Config{})
	s.Start()
	start := time.Now()
	err := s.Wait()
	elapsed := time.Since(start)

	assert.Nil(err)
	assert.True(900*time.Millisecond < elapsed)
}

func TestRegisterMiddlewares_registersExpectedMiddlewareCount(t *testing.T) {
	assert := assert.New(t)

	ms := mockEchoServer{}

	c := registerMiddlewares(&ms, 1)
	defer func() {
		c <- true
	}()

	assert.Equal(7, len(ms.middlewares))
}

func setupMockServer() *mockEchoServer {
	server := &mockEchoServer{
		groups: []groupData{},
	}

	creationFunc = func() echoServer {
		return server
	}

	return server
}

func resetCreatorFunc() {
	creationFunc = createEchoServerWrapper
}

func (m *mockRoute) Method() string {
	return m.method
}

func (m *mockRoute) Authorized() bool {
	return m.authorized
}

func (m *mockRoute) Handler() echo.HandlerFunc {
	return defaultHandler
}

func (m *mockRoute) Path() string {
	m.pathCalled++
	return m.path
}

func (m *mockEchoServer) Group(prefix string, middlewares ...echo.MiddlewareFunc) echoRouter {
	m.groupCalled++

	gd := groupData{
		prefix:      prefix,
		group:       &mockEchoRouter{},
		middlewares: len(middlewares),
	}

	m.groups = append(m.groups, gd)
	return gd.group
}

func (m *mockEchoServer) Start(address string) error {
	m.startCalled++
	m.address = address
	time.Sleep(m.sleep)
	return m.err
}

func (m *mockEchoServer) Shutdown(ctx context.Context) error {
	m.shutdownCalled++
	return m.err
}

func (m *mockEchoRouter) GET(path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *echo.Route {
	m.getCalled++
	m.path = path
	return nil
}

func (m *mockEchoRouter) POST(path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *echo.Route {
	m.postCalled++
	m.path = path
	return nil
}

func (m *mockEchoRouter) DELETE(path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *echo.Route {
	m.deleteCalled++
	m.path = path
	return nil
}

func (m *mockEchoRouter) PATCH(path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *echo.Route {
	m.patchCalled++
	m.path = path
	return nil
}

func (m *mockEchoRouter) Use(middlewares ...echo.MiddlewareFunc) {
	m.useCalled++
	m.middlewares = append(m.middlewares, middlewares...)
}
