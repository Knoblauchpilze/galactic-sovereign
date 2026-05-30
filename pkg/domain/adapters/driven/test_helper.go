package driven

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/stretchr/testify/require"
)

var (
	someTime = time.Date(2024, 11, 29, 17, 53, 29, 0, time.UTC)

	dbTestConfig = postgresql.NewConfigForLocalhost("db_galactic_sovereign", "galactic_sovereign_manager", "manager_password")
)

func newTestConnection(t *testing.T) db.Connection {
	conn, err := db.New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
}
