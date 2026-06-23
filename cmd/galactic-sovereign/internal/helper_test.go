package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testServerHost = "localhost"
)

var (
	testServerConfig = server.Config{
		BasePath:        "/v1/galactic-sovereign",
		Port:            uint16(60010),
		ShutdownTimeout: 500 * time.Millisecond,
	}

	oberonUniverseId = uuid.MustParse("9682f17b-f5f0-4eda-a747-2537d2151837")
	metalMineId      = uuid.MustParse("d176e82d-f2ca-4611-996b-c4804096caef")
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

// urlFor builds a URL under the test server's base path.
// Segments are joined with '/' and appended after the base path.
func urlFor(segments ...string) string {
	path := ""
	for _, s := range segments {
		path += "/" + s
	}
	return fmt.Sprintf(
		"http://%s:%d%s%s",
		testServerHost,
		testServerConfig.Port,
		testServerConfig.BasePath,
		path,
	)
}

func asyncStartServer(t *testing.T, s server.Server) <-chan struct{} {
	t.Helper()

	done := make(chan struct{})
	go func() {
		defer close(done)
		err := s.Start()
		assert.NoError(t, err, "Actual err: %v", err)
	}()

	const startupDelay = 50 * time.Millisecond
	time.Sleep(startupDelay)

	return done
}

func doGet[T any](t *testing.T, url string) T {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err, "GET %s: %v", url, err)
	defer resp.Body.Close() //nolint:errcheck
	require.Equal(t, http.StatusOK, resp.StatusCode, "GET %s returned %d", url, resp.StatusCode)

	return decodeResponseBody[T](t, resp.Body)
}

func doPost[T any](t *testing.T, url string, body any) T {
	t.Helper()

	raw, err := json.Marshal(body)
	require.NoError(t, err, "Actual err: %v", err)

	resp, err := http.Post(url, "application/json", bytes.NewReader(raw)) //nolint:noctx
	require.NoError(t, err, "POST %s: %v", url, err)
	defer resp.Body.Close() //nolint:errcheck
	require.Equal(t, http.StatusCreated, resp.StatusCode, "POST %s returned %d", url, resp.StatusCode)

	return decodeResponseBody[T](t, resp.Body)
}

func doDelete(t *testing.T, url string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err, "Actual err: %v", err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "DELETE %s: %v", url, err)
	defer resp.Body.Close() //nolint:errcheck
	require.Equal(t, http.StatusNoContent, resp.StatusCode, "DELETE %s returned %d", url, resp.StatusCode)
}

func decodeResponseBody[T any](t *testing.T, body io.ReadCloser) T {
	t.Helper()

	raw, err := io.ReadAll(body)
	require.NoError(t, err, "Actual err: %v", err)

	var envelope rest.ResponseEnvelope[T]
	err = json.Unmarshal(raw, &envelope)
	require.NoError(t, err, "Actual err: %v", err)

	return envelope.Details
}
