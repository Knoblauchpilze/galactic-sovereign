package internal

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestUnit_DefaultConfig(t *testing.T) {
	t.Run("defines correct rest configuration", func(t *testing.T) {
		config := DefaultConfig()

		expected := server.Config{
			BasePath:        "/v1/galactic-sovereign",
			ShutdownTimeout: 5 * time.Second,
		}
		assert.Equal(t, expected, config.Server)
	})

	t.Run("defines correct database connection", func(t *testing.T) {
		config := DefaultConfig()

		expected := postgresql.Config{
			Host:           "172.17.0.1",
			Port:           5432,
			Database:       "db_galactic_sovereign",
			User:           "galactic_sovereign_manager",
			Password:       "comes-from-the-environment",
			ConnectTimeout: 5 * time.Second,
		}
		assert.Equal(t, expected, config.Database)
	})

	t.Run("defines correct server port", func(t *testing.T) {
		config := DefaultConfig()

		assert.Equal(t, uint16(80), config.Port)
	})
}
