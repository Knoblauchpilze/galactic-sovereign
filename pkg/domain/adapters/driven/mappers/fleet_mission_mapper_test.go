package mappers

import (
	"fmt"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			var actual DbFleetMission

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

			fleet := generateSampleFleet(mission, models.StatusOutbound)
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

		fleet := generateSampleFleet(models.MissionColonize, models.StatusOutbound)

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
