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

func TestUnit_DbFleetStatus_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected models.FleetStatus
	}{
		{
			name:     "outbound string",
			input:    "OUTBOUND",
			expected: models.StatusOutbound,
		},
		{
			name:     "returning string",
			input:    "RETURNING",
			expected: models.StatusReturning,
		},
		{
			name:     "outbound bytes",
			input:    []byte("OUTBOUND"),
			expected: models.StatusOutbound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual DbFleetStatus

			err := actual.Scan(test.input)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, test.expected, actual.ToDomain())
		})
	}
}

func TestUnit_DbFleetStatus_ScanRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "unknown string",
			input: "INBOUND",
		},
		{
			name:  "unknown bytes",
			input: []byte("INBOUND"),
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
			var actual DbFleetStatus

			err := actual.Scan(test.input)

			assert.ErrorIs(t, err, domainerrors.ErrUnsupportedDatabaseEnumValue, "Actual err: %v", err)
		})
	}
}

func TestIT_DbFleetStatus(t *testing.T) {
	statuses := []models.FleetStatus{
		models.StatusOutbound,
		models.StatusReturning,
	}

	for _, status := range statuses {
		name := fmt.Sprintf("reads fleet status %s", status)
		t.Run(name, func(t *testing.T) {
			conn := newTestConnection(t)

			fleet := generateSampleFleet(models.MissionColonize, status)
			insertFleetDependencies(t, conn, fleet)
			insertTestFleet(t, conn, fleet)

			row, err := db.QueryOne[DbFleetStatus](
				t.Context(),
				conn,
				`SELECT status FROM fleet_flight WHERE fleet = $1`,
				fleet.Id,
			)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, status, row.ToDomain())
		})
	}

	t.Run("rejects unsupported values", func(t *testing.T) {
		conn := newTestConnection(t)

		fleet := generateSampleFleet(models.MissionColonize, models.StatusReturning)

		insertFleetDependencies(t, conn, fleet)
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
			"INBOUND",
		)

		dbErr, ok := db.AsDatabaseError(err)
		require.True(t, ok, "Actual err: %v", err)

		assert.Equal(t, db.ErrCheckConstraintViolation, dbErr.Code)
		assert.Equal(t, "fleet_flight", dbErr.Table)
		assert.Equal(t, "fleet_flight_status_check", dbErr.Constraint)
	})
}
