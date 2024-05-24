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

func createErrorHandlerFunc() (echo.HTTPErrorHandler, *bool, *error) {
	called := false
	var reportedErr error

	handler := func(err error, c echo.Context) {
		called = true
		reportedErr = err
	}

	return handler, &called, &reportedErr
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

func generateTestEchoContextWithErrorHandler(handler echo.HTTPErrorHandler) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.HTTPErrorHandler = handler

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)

	return ctx, rw
}

func generateTestEchoContextFromRequest(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)

	return ctx, rw
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
