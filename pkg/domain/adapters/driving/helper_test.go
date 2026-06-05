package driving

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
	someTime   = time.Date(2026, 6, 4, 21, 52, 44, 0, time.UTC)
	sampleUuid = uuid.New()
)

func generateTestRequestWithJsonBody[T any](
	t *testing.T,
	method string,
	data T,
) *http.Request {
	req := httptest.NewRequest(method, "/", encodeBody(t, data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func generateTestEchoContextFromRequest(
	t *testing.T,
	req *http.Request,
) (*echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)
	return ctx, rw
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
