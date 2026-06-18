package drivenadapters

import (
	"context"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_DatabaseChecker_Ping(t *testing.T) {
	t.Run("returns no error when connection is healthy", func(t *testing.T) {
		conn := newTestConnection(t)
		checker := NewDatabaseChecker(conn)

		err := checker.Ping(context.Background())
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when connection is not healthy", func(t *testing.T) {
		conn := newTestConnection(t)
		conn.Close(context.Background())

		checker := NewDatabaseChecker(conn)

		err := checker.Ping(context.Background())

		assert.Equal(t, db.ErrNotConnected, err, "Actual err: %v", err)
	})
}
