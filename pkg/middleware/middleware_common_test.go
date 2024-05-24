package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var defaultApiKey1 = uuid.MustParse("f847c203-1c56-43ad-9ac1-46f27d650917")
var defaultApiKey2 = uuid.MustParse("297d3309-d88b-4b83-8d82-9c6aae8a9d7a")

var errDefault = fmt.Errorf("some error")

func createHandlerFuncWithCalledBoolean() (echo.HandlerFunc, *bool) {
	called := false
	call := func(c echo.Context) error {
		called = true
		return nil
	}
	return call, &called
}

func createHandlerFuncReturning(err error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return err
	}
}

func createHandlerFuncReturningCode(code int) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(code, "")
	}
}

func generateTestEchoContextWithApiKey() (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(apiKeyHeaderKey, defaultApiKey1.String())
	return generateTestEchoContextFromRequest(req)
}

func generateTestEchoContext() (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	return generateTestEchoContextFromRequest(req)
}

func generateTestEchoContextFromRequest(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)
	return ctx, rw
}

type mockEchoContext struct {
	// https://www.myhatchpad.com/insight/mocking-techniques-for-go/
	// See Embedding Interfaces in the page.
	echo.Context

	request        *http.Request
	response       *echo.Response
	logger         echo.Logger
	values         map[string]interface{}
	loggerChanged  bool
	requestChanged bool
	reportedError  error

	reportedCode int
	jsonContent  interface{}
}

type mockEchoLogger struct {
	// https://stackoverflow.com/questions/34754210/why-does-embedding-an-interface-in-a-struct-cause-the-interface-methodset-to-be
	echo.Logger
}

type mockResponseWriter struct {
	headerCalled int
	header       http.Header

	writeCalled int
	data        []byte
	written     int
	writeErr    error

	writeHeaderCalled int
	code              int
}

func newMockEchoContext(code int) *mockEchoContext {
	return &mockEchoContext{
		request: &http.Request{
			Method: "GET",
		},
		response: &echo.Response{
			Status: code,
			Writer: &mockResponseWriter{},
		},
		logger: &mockEchoLogger{},
		values: map[string]interface{}{},
	}
}

func (m *mockEchoContext) Request() *http.Request {
	return m.request
}

func (m *mockEchoContext) SetRequest(req *http.Request) {
	m.requestChanged = true
	m.request = req
}

func (m *mockEchoContext) Response() *echo.Response {
	return m.response
}

func (m *mockEchoContext) Set(key string, val interface{}) {
	m.values[key] = val
}

func (m *mockEchoContext) JSON(code int, i interface{}) error {
	m.reportedCode = code
	m.jsonContent = i
	return nil
}

func (m *mockEchoContext) Error(err error) {
	m.reportedError = err
}

func (m *mockEchoContext) Logger() echo.Logger {
	return m.logger
}

func (m *mockEchoContext) SetLogger(l echo.Logger) {
	m.loggerChanged = true
}

func (m *mockEchoLogger) Warnf(format string, args ...interface{}) {}

func (m *mockEchoLogger) Errorf(format string, args ...interface{}) {}

func (m *mockResponseWriter) Header() http.Header {
	m.headerCalled++
	return m.header
}

func (m *mockResponseWriter) Write(out []byte) (int, error) {
	m.writeCalled++
	m.data = out
	return m.written, m.writeErr
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.writeHeaderCalled++
	m.code = statusCode
}
