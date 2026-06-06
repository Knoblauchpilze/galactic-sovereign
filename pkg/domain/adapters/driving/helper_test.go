package drivingadapters

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

var (
	someTime      = time.Date(2026, 6, 4, 21, 52, 44, 0, time.UTC)
	someOtherTime = time.Date(2026, 6, 5, 18, 13, 10, 0, time.UTC)
	sampleUuid    = uuid.New()
)

func generateTestRequest(t *testing.T, method string) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, "/", nil)
	return req
}

func generateTestRequestWithJsonBody[T any](
	t *testing.T,
	method string,
	data T,
) *http.Request {
	req := httptest.NewRequest(method, "/", encodeBody(t, data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func generateTestContextFromRequest(
	t *testing.T,
	req *http.Request,
	modifiers ...func(*testing.T, *echo.Context),
) (*echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)

	for _, modifier := range modifiers {
		modifier(t, ctx)
	}

	return ctx, rw
}

func addIdPathParam(t *testing.T, c *echo.Context) {
	t.Helper()

	c.SetPathValues([]echo.PathValue{{Name: "id", Value: sampleUuid.String()}})
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
