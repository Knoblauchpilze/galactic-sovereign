package middleware

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
	request       *http.Request
	response      *echo.Response
	logger        echo.Logger
	reportedError error
}

type mockEchoLogger struct{}

func newMockEchoContext(req *http.Request, res *echo.Response) *mockEchoContext {
	return &mockEchoContext{
		request:  req,
		response: res,
		logger:   &mockEchoLogger{},
	}
}

func (m *mockEchoContext) Request() *http.Request { return m.request }

func (m *mockEchoContext) SetRequest(r *http.Request) {}

func (m *mockEchoContext) SetResponse(r *echo.Response) {}

func (m *mockEchoContext) Response() *echo.Response { return m.response }

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

func (m *mockEchoContext) Error(err error) {
	m.reportedError = err
}

func (m *mockEchoContext) Handler() echo.HandlerFunc { return nil }

func (m *mockEchoContext) SetHandler(h echo.HandlerFunc) {}

func (m *mockEchoContext) Logger() echo.Logger { return m.logger }

func (m *mockEchoContext) SetLogger(l echo.Logger) {}

func (m *mockEchoContext) Echo() *echo.Echo { return nil }

func (m *mockEchoContext) Reset(r *http.Request, w http.ResponseWriter) {}

func (m *mockEchoLogger) Output() io.Writer { return nil }

func (m *mockEchoLogger) SetOutput(w io.Writer) {}

func (m *mockEchoLogger) Prefix() string { return "" }

func (m *mockEchoLogger) SetPrefix(p string) {}

func (m *mockEchoLogger) Level() log.Lvl { return log.DEBUG }

func (m *mockEchoLogger) SetLevel(v log.Lvl) {}

func (m *mockEchoLogger) SetHeader(h string) {}

func (m *mockEchoLogger) Print(i ...interface{}) {}

func (m *mockEchoLogger) Printf(format string, args ...interface{}) {}

func (m *mockEchoLogger) Printj(j log.JSON) {}

func (m *mockEchoLogger) Debug(i ...interface{}) {}

func (m *mockEchoLogger) Debugf(format string, args ...interface{}) {}

func (m *mockEchoLogger) Debugj(j log.JSON) {}

func (m *mockEchoLogger) Info(i ...interface{}) {}

func (m *mockEchoLogger) Infof(format string, args ...interface{}) {}

func (m *mockEchoLogger) Infoj(j log.JSON) {}

func (m *mockEchoLogger) Warn(i ...interface{}) {}

func (m *mockEchoLogger) Warnf(format string, args ...interface{}) {}

func (m *mockEchoLogger) Warnj(j log.JSON) {}

func (m *mockEchoLogger) Error(i ...interface{}) {}

func (m *mockEchoLogger) Errorf(format string, args ...interface{}) {}

func (m *mockEchoLogger) Errorj(j log.JSON) {}

func (m *mockEchoLogger) Fatal(i ...interface{}) {}

func (m *mockEchoLogger) Fatalj(j log.JSON) {}

func (m *mockEchoLogger) Fatalf(format string, args ...interface{}) {}

func (m *mockEchoLogger) Panic(i ...interface{}) {}

func (m *mockEchoLogger) Panicj(j log.JSON) {}

func (m *mockEchoLogger) Panicf(format string, args ...interface{}) {}
