package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/db/postgresql"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("db_galactic_sovereign", "galactic_sovereign_manager", "JT!vF37s7vj#^%eZjHTSzKs49HCaz")

func newTestConnection(t *testing.T) db.Connection {
	conn, err := db.New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
}

func insertTestBuildingAction(t *testing.T, conn db.Connection) (persistence.BuildingAction, persistence.Planet) {
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building := insertTestBuilding(t, conn)

	action := persistence.BuildingAction{
		Id:           uuid.New(),
		Planet:       planet.Id,
		Building:     building.Id,
		CurrentLevel: 4,
		DesiredLevel: 5,
		CreatedAt:    someTime,
		CompletedAt:  someTime.Add(1*time.Hour + 2*time.Minute),
	}

	sqlQuery := `INSERT INTO building_action (id, planet, building, current_level, desired_level, created_at, completed_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		action.Id,
		action.Planet,
		action.Building,
		action.CurrentLevel,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
	)
	require.Nil(t, err)

	return action, planet
}
