package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	testServerHost = "localhost"
)

var (
	oberonUniverseId = uuid.MustParse("9682f17b-f5f0-4eda-a747-2537d2151837")
	metalMineId      = uuid.MustParse("d176e82d-f2ca-4611-996b-c4804096caef")
	shipyardId       = uuid.MustParse("58d75842-6dc0-4ac0-b36d-55f91b8d060d")
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

// urlFor builds a URL under the test server's base path.
// Segments are joined with '/' and appended after the base path.
func urlFor(conf server.Config, segments ...string) string {
	path := ""
	for _, s := range segments {
		path += "/" + s
	}
	return fmt.Sprintf(
		"http://%s:%d%s%s",
		testServerHost,
		conf.Port,
		conf.BasePath,
		path,
	)
}

func newTestServerConfig() server.Config {
	return server.Config{
		BasePath:        "/v1/galactic-sovereign",
		ShutdownTimeout: 500 * time.Millisecond,
	}
}

// asyncStartServer binds the server on an OS-assigned port, serves until test
// cleanup and returns the config updated with the effective port.
func asyncStartServer(t *testing.T, s HttpServer, conf server.Config) server.Config {
	t.Helper()

	listener, err := s.Bind(0)
	require.NoError(t, err, "Actual err: %v", err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- s.Serve(ctx, listener)
	}()

	t.Cleanup(func() {
		cancel()
		err := <-done
		require.NoError(t, err, "Actual err: %v", err)
	})

	conf.Port = uint16(listener.Addr().(*net.TCPAddr).Port)
	return conf
}

func doGet[T any](t *testing.T, url string) T {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err, "GET %s: %v", url, err)
	defer resp.Body.Close() // nolint:errcheck
	require.Equal(t, http.StatusOK, resp.StatusCode, "GET %s returned %d", url, resp.StatusCode)

	return decodeResponseBody[T](t, resp.Body)
}

func doPost[T any](t *testing.T, url string, body any) T {
	t.Helper()

	var payload io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		require.NoError(t, err, "Actual err: %v", err)

		payload = bytes.NewReader(raw)
	}

	resp, err := http.Post(url, "application/json", payload) // nolint:noctx
	require.NoError(t, err, "POST %s: %v", url, err)
	defer resp.Body.Close() // nolint:errcheck
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
	defer resp.Body.Close() // nolint:errcheck
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

// addPlanetResources credits the given resource amount to the planet, allowing
// tests to afford actions (e.g. ships) that starting resources can't cover.
func addPlanetResources(
	t *testing.T,
	conn database.Connection,
	planet uuid.UUID,
	resource uuid.UUID,
	amount int,
) {
	t.Helper()

	_, err := conn.Exec(
		t.Context(),
		`UPDATE planet_resource SET amount = amount + $1 WHERE planet = $2 AND resource = $3`,
		amount,
		planet,
		resource,
	)
	require.NoError(t, err, "Actual err: %v", err)
}

// bumpPlanetBuilding bumpes the given building to the desired level on the
// planet, allowing tests to perform actions (e.g. ship creation) that the
// starting planet would not allow.
func bumpPlanetBuilding(
	t *testing.T,
	conn database.Connection,
	planet uuid.UUID,
	building uuid.UUID,
	level int,
) {
	t.Helper()

	_, err := conn.Exec(
		t.Context(),
		`UPDATE planet_building SET level = $1 WHERE planet = $2 AND building = $3`,
		level,
		planet,
		building,
	)
	require.NoError(t, err, "Actual err: %v", err)
}
