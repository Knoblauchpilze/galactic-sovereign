package internal

import (
	"os"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
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
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
