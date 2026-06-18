package drivenadapters

import (
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_DatabaseChecker_Ping(t *testing.T) {
	t.Run("returns no error when connection is healthy", func(t *testing.T) {
		conn := newTestConnection(t)
		checker := NewDatabaseChecker(conn)

		err := checker.Ping(t.Context())
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when connection is not healthy", func(t *testing.T) {
		conn := newTestConnection(t)
		conn.Close(t.Context())

		checker := NewDatabaseChecker(conn)

		err := checker.Ping(t.Context())

		assert.Equal(t, db.ErrNotConnected, err, "Actual err: %v", err)
	})
}
