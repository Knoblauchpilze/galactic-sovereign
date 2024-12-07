package repositories

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/db/postgresql"
	"github.com/stretchr/testify/require"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("db_galactic_sovereign", "galactic_sovereign_manager", "JT!vF37s7vj#^%eZjHTSzKs49HCaz")

func newTestConnection(t *testing.T) db.Connection {
	conn, err := db.New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
}
