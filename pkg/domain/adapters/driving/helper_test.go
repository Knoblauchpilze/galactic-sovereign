package drivingadapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	someTime         = time.Date(2026, time.June, 4, 21, 52, 44, 0, time.UTC)
	someOtherTime    = time.Date(2026, time.June, 5, 18, 13, 10, 0, time.UTC)
	sampleUuid       = uuid.New()
	sampleResourceId = uuid.New()
)

func generateTestRequest(
	t *testing.T,
	method string,
	modifiers ...func(*testing.T, *http.Request),
) *http.Request {
	t.Helper()

	ctx := rest.WithContextLogger(t.Context(), slog.Default())
	req := httptest.NewRequestWithContext(ctx, method, "/", nil)

	for _, modifier := range modifiers {
		modifier(t, req)
	}

	return req
}

func addSampleUuidPathParam(t *testing.T, req *http.Request) {
	t.Helper()

	req.URL.Path = fmt.Sprintf("/%s", sampleUuid)
}

func addInvalidUuidPathParam(t *testing.T, req *http.Request) {
	t.Helper()

	req.URL.Path = "/not-a-uuid"
}

func addRequestPath(t *testing.T, req *http.Request, path string, args ...any) {
	t.Helper()

	req.URL.Path = fmt.Sprintf(path, args...)
}

func generateTestRequestWithJsonBody[T any](
	t *testing.T,
	method string,
	data T,
) *http.Request {
	ctx := rest.WithContextLogger(t.Context(), slog.Default())
	req := httptest.NewRequestWithContext(ctx, method, "/", encodeBody(t, data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func createTestGinRouter(
	t *testing.T,
	method string,
	path string,
	handler gin.HandlerFunc,
	middlewares ...gin.HandlerFunc,
) *gin.Engine {
	t.Helper()

	r := gin.New()

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	r.Handle(method, path, handler)

	return r
}

func decodeResponseBody[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()

	var responseBody T

	rawBody, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err, "Actual err: %v", err)

	err = json.Unmarshal(rawBody, &responseBody)
	require.NoError(t, err, "Actual err: %v", err)

	return responseBody
}

func encodeBody[T any](t *testing.T, data T) io.Reader {
	t.Helper()

	out, err := json.Marshal(data)
	require.NoError(t, err, "Actual err: %v", err)

	return bytes.NewReader(out)
}

func ptrFor[T any](value T) *T {
	return &value
}
