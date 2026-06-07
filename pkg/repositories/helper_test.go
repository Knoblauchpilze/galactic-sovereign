package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func insertTestPlayer(t *testing.T, conn db.Connection) uuid.UUID {
	someTime := time.Date(2024, 11, 29, 17, 56, 02, 0, time.UTC)

	playerId := uuid.New()

	sqlQuery := `INSERT INTO player (id, api_user, universe, name, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		playerId,
		uuid.New(),
		oberonUniverseId,
		fmt.Sprintf("my-player-%s", playerId.String()),
		someTime,
	)
	require.Nil(t, err)

	return playerId
}

func insertTestPlanet(t *testing.T, conn db.Connection, player uuid.UUID, homeworld bool) uuid.UUID {
	someTime := time.Date(2024, 11, 30, 11, 31, 58, 0, time.UTC)

	planetId := uuid.New()

	sqlQuery := `INSERT INTO planet (id, player, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planetId,
		player,
		fmt.Sprintf("my-planet-%s", planetId.String()),
		someTime,
		someTime,
	)
	require.Nil(t, err)

	if homeworld {
		sqlQuery := `INSERT INTO homeworld (player, planet) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), sqlQuery, player, planetId)
		require.Nil(t, err)
	}

	return planetId
}

func insertTestPlanetForPlayer(t *testing.T, conn db.Connection) (uuid.UUID, uuid.UUID) {
	playerId := insertTestPlayer(t, conn)
	planetId := insertTestPlanet(t, conn, playerId, false)
	return planetId, playerId
}
