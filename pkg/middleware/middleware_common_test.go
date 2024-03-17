package middleware

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

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

type mockEchoContext struct{}

func (m *mockEchoContext) Request() *http.Request { return nil }

func (m *mockEchoContext) SetRequest(r *http.Request) {}

func (m *mockEchoContext) SetResponse(r *echo.Response) {}

func (m *mockEchoContext) Response() *echo.Response { return nil }

func (m *mockEchoContext) IsTLS() bool { return false }

func (m *mockEchoContext) IsWebSocket() bool { return false }

func (m *mockEchoContext) Scheme() string { return "" }

func (m *mockEchoContext) RealIP() string { return "" }

func (m *mockEchoContext) Path() string { return "" }

func (m *mockEchoContext) SetPath(p string) {}

func (m *mockEchoContext) Param(name string) string { return "" }

func (m *mockEchoContext) ParamNames() []string { return nil }

func (m *mockEchoContext) SetParamNames(names ...string) {}

func (m *mockEchoContext) ParamValues() []string { return nil }

func (m *mockEchoContext) SetParamValues(values ...string) {}

func (m *mockEchoContext) QueryParam(name string) string { return "" }

func (m *mockEchoContext) QueryParams() url.Values { return nil }

func (m *mockEchoContext) QueryString() string { return "" }

func (m *mockEchoContext) FormValue(name string) string { return "" }

func (m *mockEchoContext) FormParams() (url.Values, error) { return nil, nil }

func (m *mockEchoContext) FormFile(name string) (*multipart.FileHeader, error) { return nil, nil }

func (m *mockEchoContext) MultipartForm() (*multipart.Form, error) { return nil, nil }

func (m *mockEchoContext) Cookie(name string) (*http.Cookie, error) { return nil, nil }

func (m *mockEchoContext) SetCookie(cookie *http.Cookie) {}

func (m *mockEchoContext) Cookies() []*http.Cookie { return nil }

func (m *mockEchoContext) Get(key string) interface{} { return nil }

func (m *mockEchoContext) Set(key string, val interface{}) {}

func (m *mockEchoContext) Bind(i interface{}) error { return nil }

func (m *mockEchoContext) Validate(i interface{}) error { return nil }

func (m *mockEchoContext) Render(code int, name string, data interface{}) error { return nil }

func (m *mockEchoContext) HTML(code int, html string) error { return nil }

func (m *mockEchoContext) HTMLBlob(code int, b []byte) error { return nil }

func (m *mockEchoContext) String(code int, s string) error { return nil }

func (m *mockEchoContext) JSON(code int, i interface{}) error { return nil }

func (m *mockEchoContext) JSONPretty(code int, i interface{}, indent string) error { return nil }

func (m *mockEchoContext) JSONBlob(code int, b []byte) error { return nil }

func (m *mockEchoContext) JSONP(code int, callback string, i interface{}) error { return nil }

func (m *mockEchoContext) JSONPBlob(code int, callback string, b []byte) error { return nil }

func (m *mockEchoContext) XML(code int, i interface{}) error { return nil }

func (m *mockEchoContext) XMLPretty(code int, i interface{}, indent string) error { return nil }

func (m *mockEchoContext) XMLBlob(code int, b []byte) error { return nil }

func (m *mockEchoContext) Blob(code int, contentType string, b []byte) error { return nil }

func (m *mockEchoContext) Stream(code int, contentType string, r io.Reader) error { return nil }

func (m *mockEchoContext) File(file string) error { return nil }

func (m *mockEchoContext) Attachment(file string, name string) error { return nil }

func (m *mockEchoContext) Inline(file string, name string) error { return nil }

func (m *mockEchoContext) NoContent(code int) error { return nil }

func (m *mockEchoContext) Redirect(code int, url string) error { return nil }

func (m *mockEchoContext) Error(err error) {}

func (m *mockEchoContext) Handler() echo.HandlerFunc { return nil }

func (m *mockEchoContext) SetHandler(h echo.HandlerFunc) {}

func (m *mockEchoContext) Logger() echo.Logger { return nil }

func (m *mockEchoContext) SetLogger(l echo.Logger) {}

func (m *mockEchoContext) Echo() *echo.Echo { return nil }

func (m *mockEchoContext) Reset(r *http.Request, w http.ResponseWriter) {}
