package drivenadapters

import (
	"context"
	"math"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/stretchr/testify/require"
)

var (
	someTime      = time.Date(2024, 11, 29, 17, 53, 29, 0, time.UTC)
	someOtherTime = time.Date(2026, 6, 1, 8, 20, 15, 0, time.UTC)

	dbTestConfig = postgresql.NewConfigForLocalhost("db_galactic_sovereign", "galactic_sovereign_manager", "manager_password")
)

func newTestConnection(t *testing.T) db.Connection {
	conn, err := db.New(context.Background(), dbTestConfig)
	require.NoError(t, err, "Actual err: %v", err)
	return conn
}

func randFloat(precision int) float64 {
	rounder := math.Pow(10, float64(precision))
	return math.Round(rand.Float64()*rounder) / rounder
}
