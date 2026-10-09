package mappers

import (
	"fmt"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	t1 = time.Date(2026, time.October, 3, 11, 9, 40, 0, time.UTC)
)

func TestUnit_DbFleetMission_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected models.FleetMission
	}{
		{
			name:     "transport string",
			input:    "TRANSPORT",
			expected: models.MissionTransport,
		},
		{
			name:     "colonize string",
			input:    "COLONIZE",
			expected: models.MissionColonize,
		},
		{
			name:     "colonize bytes",
			input:    []byte("COLONIZE"),
			expected: models.MissionColonize,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual DbFleetMission

			err := actual.Scan(test.input)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, test.expected, actual.ToDomain())
		})
	}
}

func TestUnit_DbFleetMission_ScanRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "unknown string",
			input: "DEPOSIT",
		},
		{
			name:  "unknown bytes",
			input: []byte("DEPOSIT"),
		},
		{
			name:  "nil",
			input: nil,
		},
		{
			name:  "unsupported type",
			input: 42,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual DbScalingMode

			err := actual.Scan(test.input)

			assert.ErrorIs(t, err, domainerrors.ErrUnsupportedDatabaseEnumValue, "Actual err: %v", err)
		})
	}
}

func TestIT_DbFleetMission(t *testing.T) {
	missions := []models.FleetMission{
		models.MissionTransport,
		models.MissionColonize,
	}

	for _, mission := range missions {
		name := fmt.Sprintf("reads fleet mission %s", mission)
		t.Run(name, func(t *testing.T) {
			conn := newTestConnection(t)

			fleet := generateSampleFleet(mission)
			insertFleetDependencies(t, conn, fleet)
			insertTestFleet(t, conn, fleet)

			row, err := db.QueryOne[DbFleetMission](
				t.Context(),
				conn,
				`SELECT mission FROM fleet WHERE id = $1`,
				fleet.Id,
			)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, mission, row.ToDomain())
		})
	}

	t.Run("rejects unsupported values", func(t *testing.T) {
		conn := newTestConnection(t)

		fleet := generateSampleFleet(models.MissionColonize)

		_, err := conn.Exec(
			t.Context(),
			`INSERT INTO fleet (id, player, source, mission, created_at, updated_at, version)
					VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			fleet.Id,
			fleet.Player,
			fleet.Source,
			"DEPOSIT",
			fleet.CreatedAt,
			fleet.UpdatedAt,
			fleet.Version,
		)

		dbErr, ok := db.AsDatabaseError(err)
		require.True(t, ok, "Actual err: %v", err)

		assert.Equal(t, db.ErrCheckConstraintViolation, dbErr.Code)
		assert.Equal(t, "fleet", dbErr.Table)
		assert.Equal(t, "fleet_mission_check", dbErr.Constraint)
	})
}

func generateSampleFleet(mission models.FleetMission) models.Fleet {
	return models.Fleet{
		Id:          uuid.New(),
		Player:      uuid.New(),
		Source:      uuid.New(),
		Mission:     mission,
		Destination: models.FleetDestination{},
		Ships:       []models.FleetShip{},
		CreatedAt:   t1,
		ArrivalAt:   t1.Add(2 * time.Minute),
		ReturnAt:    t1.Add(4 * time.Minute),
		UpdatedAt:   t1,
		Version:     0,
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
