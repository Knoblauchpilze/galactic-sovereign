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

func TestUnit_DbFleetEventType_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected models.FleetEventType
	}{
		{
			name:     "arrival string",
			input:    "ARRIVAL",
			expected: models.ArrivalEvent,
		},
		{
			name:     "return string",
			input:    "RETURN",
			expected: models.ReturnEvent,
		},
		{
			name:     "arrival bytes",
			input:    []byte("ARRIVAL"),
			expected: models.ArrivalEvent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual DbFleetEventType

			err := actual.Scan(test.input)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, test.expected, actual.ToDomain())
		})
	}
}

func TestUnit_DbFleetEventType_ScanRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "unknown string",
			input: "DESTRUCTION",
		},
		{
			name:  "unknown bytes",
			input: []byte("DESTRUCTION"),
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
			var actual DbFleetEventType

			err := actual.Scan(test.input)

			assert.ErrorIs(t, err, domainerrors.ErrUnsupportedDatabaseEnumValue, "Actual err: %v", err)
		})
	}
}

func TestIT_DbFleetEventType(t *testing.T) {
	events := []models.FleetEventType{
		models.ArrivalEvent,
		models.ReturnEvent,
	}

	for _, event := range events {
		name := fmt.Sprintf("reads fleet event type %s", event)
		t.Run(name, func(t *testing.T) {
			conn := newTestConnection(t)

			fleet := generateSampleFleet(models.MissionColonize, models.StatusOutbound)
			insertFleetDependencies(t, conn, fleet)
			insertTestFleet(t, conn, fleet)

			_, err := conn.Exec(
				t.Context(),
				`INSERT INTO fleet_event
						(fleet, next_event_at, event_type)
						VALUES ($1, $2, $3)`,
				fleet.Id,
				t1,
				event,
			)
			require.NoError(t, err, "Actual err: %v", err)

			row, err := db.QueryOne[DbFleetEventType](
				t.Context(),
				conn,
				`SELECT event_type
					 FROM fleet_event
					 WHERE fleet = $1`,
				fleet.Id,
			)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, event, row.ToDomain())
		})
	}

	t.Run("rejects unsupported values", func(t *testing.T) {
		conn := newTestConnection(t)

		fleet := generateSampleFleet(models.MissionColonize, models.StatusOutbound)
		insertFleetDependencies(t, conn, fleet)
		insertTestFleet(t, conn, fleet)

		_, err := conn.Exec(
			t.Context(),
			`INSERT INTO fleet_event
						(fleet, next_event_at, event_type)
						VALUES ($1, $2, $3)`,
			fleet.Id,
			t1,
			"DESTRUCTION",
		)

		dbErr, ok := db.AsDatabaseError(err)
		require.True(t, ok, "Actual err: %v", err)

		assert.Equal(t, db.ErrCheckConstraintViolation, dbErr.Code)
		assert.Equal(t, "fleet_event", dbErr.Table)
		assert.Equal(t, "fleet_event_event_type_check", dbErr.Constraint)
	})
}
