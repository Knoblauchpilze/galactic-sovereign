package drivingadapters

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	sampleKey = "my-key"
)

func TestUnit_FetchIdFromQueryParam(t *testing.T) {
	t.Run("returns not found/no error when no query parameter provided", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, _ := generateTestContextFromRequest(t, req)

		exists, _, err := fetchIdFromQueryParam(sampleKey, ctx)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, exists)
	})

	t.Run("returns not found/no error when id set for a different key", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "not-the-default-key", sampleUuid.String())
		ctx, _ := generateTestContextFromRequest(t, req)

		exists, _, err := fetchIdFromQueryParam(sampleKey, ctx)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, exists)
	})

	t.Run("returns exists/error when id is set with wrong syntax", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, sampleKey, "not-a-uuid")
		ctx, _ := generateTestContextFromRequest(t, req)

		exists, _, err := fetchIdFromQueryParam(sampleKey, ctx)

		assert.True(t, exists)
		assert.True(t, uuid.IsInvalidLengthError(err), "Actual err: %v", err)
	})

	t.Run("returns exists/no error when id is set and valid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, sampleKey, sampleUuid.String())
		ctx, _ := generateTestContextFromRequest(t, req)

		exists, actual, err := fetchIdFromQueryParam(sampleKey, ctx)
		require.NoError(t, err, "Actual err: %v", err)

		assert.True(t, exists)
		assert.Equal(t, sampleUuid, actual)
	})
}

func addQueryParam(t *testing.T, req *http.Request, key string, value string) {
	t.Helper()

	q := req.URL.Query()
	q.Add(key, value)

	req.URL.RawQuery = q.Encode()
}
