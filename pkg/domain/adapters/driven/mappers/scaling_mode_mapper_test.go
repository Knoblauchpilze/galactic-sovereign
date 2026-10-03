package mappers

import (
	"fmt"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_DbScalingMode_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected models.ScalingMode
	}{
		{
			name:     "linear string",
			input:    "LINEAR",
			expected: models.LinearScaling,
		},
		{
			name:     "geometric string",
			input:    "GEOMETRIC",
			expected: models.GeometricScaling,
		},
		{
			name:     "linear bytes",
			input:    []byte("LINEAR"),
			expected: models.LinearScaling,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual DbScalingMode

			err := actual.Scan(test.input)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, test.expected, actual.ToDomain())
		})
	}
}

func TestUnit_DbScalingMode_ScanRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "unknown string",
			input: "EXPONENTIAL",
		},
		{
			name:  "unknown bytes",
			input: []byte("EXPONENTIAL"),
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

func TestIT_DbScalingMode(t *testing.T) {
	modes := []models.ScalingMode{
		models.LinearScaling,
		models.GeometricScaling,
	}

	for _, mode := range modes {
		name := fmt.Sprintf("reads scaling mode %s", mode)
		t.Run(name, func(t *testing.T) {
			conn := newTestConnection(t)

			buildingId := insertTestBuilding(t, conn)

			_, err := conn.Exec(
				t.Context(),
				`INSERT INTO building_resource_metabolization_ship_speedup
						(building, scaling, base, coefficient)
						VALUES ($1, $2, $3, $4)`,
				buildingId,
				mode,
				1.0,
				1.0,
			)
			require.NoError(t, err, "Actual err: %v", err)

			row, err := db.QueryOne[DbScalingMode](
				t.Context(),
				conn,
				`SELECT scaling
					 FROM building_resource_metabolization_ship_speedup
					 WHERE building = $1`,
				buildingId,
			)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, mode, row.ToDomain())
		})
	}

	t.Run("rejects unsupported values", func(t *testing.T) {
		conn := newTestConnection(t)

		buildingId := insertTestBuilding(t, conn)

		_, err := conn.Exec(
			t.Context(),
			`INSERT INTO building_resource_metabolization_ship_speedup
				(building, scaling, base, coefficient)
				VALUES ($1, $2, $3, $4)`,
			buildingId,
			"EXPONENTIAL",
			1.0,
			1.0,
		)

		dbErr, ok := db.AsDatabaseError(err)
		require.True(t, ok, "Actual err: %v", err)

		assert.Equal(t, db.ErrCheckConstraintViolation, dbErr.Code)
		assert.Equal(t, "building_resource_metabolization_ship_speedup", dbErr.Table)
		assert.Equal(t, "building_resource_metabolization_ship_speedup_scaling_check", dbErr.Constraint)
	})
}

func insertTestBuilding(
	t *testing.T,
	conn database.Connection,
) uuid.UUID {
	t.Helper()

	buildingId := uuid.New()

	_, err := conn.Exec(
		t.Context(),
		`INSERT INTO building (id, name) VALUES ($1, $2)`,
		buildingId,
		fmt.Sprintf("scaling-mode-test-%s", buildingId),
	)
	require.NoError(t, err, "Actual err: %v", err)

	return buildingId
}
