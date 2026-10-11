package mappers

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	integrationdb "github.com/Knoblauchpilze/galactic-sovereign/pkg/testing/integrationdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	sharedDbContainer = &integrationdb.Suite{}

	t1 = time.Date(2026, time.October, 3, 11, 9, 40, 0, time.UTC)
)

func TestMain(m *testing.M) {
	code := m.Run()
	sharedDbContainer.Teardown()
	os.Exit(code)
}

func newTestConnection(t *testing.T) database.Connection {
	t.Helper()
	return sharedDbContainer.NewTestConnection(t)
}

func generateSampleFleet(mission models.FleetMission, status models.FleetStatus) models.Fleet {
	return models.Fleet{
		Id:          uuid.New(),
		Player:      uuid.New(),
		Source:      uuid.New(),
		Mission:     mission,
		Destination: models.FleetDestination{},
		Ships:       []models.FleetShip{},
		Schedule: models.FleetFlightSchedule{
			ArrivalAt: t1.Add(2 * time.Minute),
			ReturnAt:  t1.Add(4 * time.Minute),
			Status:    status,
		},
		CreatedAt: t1,
		UpdatedAt: t1,
		Version:   0,
	}
}

func insertTestFleet(t *testing.T, conn database.Connection, fleet models.Fleet) {
	t.Helper()

	_, err := conn.Exec(
		t.Context(),
		`INSERT INTO fleet (id, player, source, mission, created_at, updated_at, version)
					VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		fleet.Id,
		fleet.Player,
		fleet.Source,
		fleet.Mission,
		fleet.CreatedAt,
		fleet.UpdatedAt,
		fleet.Version,
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = conn.Exec(
		t.Context(),
		`INSERT INTO fleet_flight (fleet, arrival_at, return_at, status)
					VALUES ($1, $2, $3, $4)`,
		fleet.Id,
		fleet.Schedule.ArrivalAt,
		fleet.Schedule.ReturnAt,
		fleet.Schedule.Status,
	)
	require.NoError(t, err, "Actual err: %v", err)
}

// insertFleetDependencies will create a fresh universe and registers the player to whom the
// fleet belongs in it. It will also registers a planet with an identifier equivalent to the
// source of the fleet. This allows to make inserting the fleet possible by fulfilling the
// foreign keys.
// Note: the target is not inserted and is considered nil.
func insertFleetDependencies(
	t *testing.T,
	conn database.Connection,
	fleet models.Fleet,
) {
	t.Helper()

	universeId := uuid.New()

	_, err := conn.Exec(
		t.Context(),
		`INSERT INTO universe (id, name) VALUES ($1, $2)`,
		universeId,
		fmt.Sprintf("test-universe-%s", universeId),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = conn.Exec(
		t.Context(),
		`INSERT INTO player (id, api_user, universe, name) VALUES ($1, $2, $3, $4)`,
		fleet.Player,
		uuid.New(),
		universeId,
		fmt.Sprintf("test-player-%s", fleet.Player),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = conn.Exec(
		t.Context(),
		`INSERT INTO planet (id, player, name, fields, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		fleet.Source,
		fleet.Player,
		fmt.Sprintf("test-planet-%s", fleet.Source),
		36,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	require.NoError(t, err, "Actual err: %v", err)
}
