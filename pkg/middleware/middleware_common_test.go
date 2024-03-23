package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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

type mockEchoContext struct {
	// https://www.myhatchpad.com/insight/mocking-techniques-for-go/
	// See Embedding Interfaces in the page.
	echo.Context

	request       *http.Request
	response      *echo.Response
	logger        echo.Logger
	values        map[string]interface{}
	loggerChanged bool
	reportedError error
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

func (m *mockEchoContext) Request() *http.Request { return m.request }

func (m *mockEchoContext) Response() *echo.Response { return m.response }

func (m *mockEchoContext) Set(key string, val interface{}) {
	m.values[key] = val
}

func (m *mockEchoContext) Error(err error) {
	m.reportedError = err
}

func (m *mockEchoContext) Logger() echo.Logger { return m.logger }

func (m *mockEchoContext) SetLogger(l echo.Logger) {
	m.loggerChanged = true
}

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
