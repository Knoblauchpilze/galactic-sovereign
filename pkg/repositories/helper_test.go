package repositories

import (
	"context"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	oberonUniverseId = uuid.MustParse("9682f17b-f5f0-4eda-a747-2537d2151837")
	dbTestConfig     = postgresql.NewConfigForLocalhost("db_galactic_sovereign", "galactic_sovereign_manager", "manager_password")
)

func newTestConnection(t *testing.T) db.Connection {
	conn, err := db.New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
}
