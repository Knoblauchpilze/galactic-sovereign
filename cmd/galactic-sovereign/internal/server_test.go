package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	integrationdb "github.com/Knoblauchpilze/galactic-sovereign/pkg/testing/integrationdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_Server_GetUniverses(t *testing.T) {
	dbContainer := integrationdb.NewDatabaseSharedContainer(t)

	conn := dbContainer.NewTestConnection(t)

	s := CreateGameServer(testServerConfig, conn, slog.Default())
	done := asyncStartServer(t, s)

	url := fmt.Sprintf("http://%s:%d/v1/galactic-sovereign/universes", testServerHost, testServerConfig.Port)
	resp, err := http.Get(url)
	require.NoError(t, err, "Actual err: %v", err)
	// nolint: errcheck
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = s.Stop()
	require.NoError(t, err, "Actual err: %v", err)
	<-done
}

// asyncStartServer starts the server in a background goroutine and waits briefly
// for it to be ready before returning. The returned channel is closed when the
// server exits; callers should read it after calling s.Stop() to catch any error
// from s.Start().
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
