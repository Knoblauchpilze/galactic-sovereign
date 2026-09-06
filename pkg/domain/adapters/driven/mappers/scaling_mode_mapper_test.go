package mappers

import (
	"fmt"
	"os"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	integrationdb "github.com/Knoblauchpilze/galactic-sovereign/pkg/testing/integrationdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sharedDbContainer = &integrationdb.Suite{}

func TestMain(m *testing.M) {
	code := m.Run()
	sharedDbContainer.Teardown()
	os.Exit(code)
}

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

			require.NoError(t, err)
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
			actual := DbScalingMode(models.LinearScaling)

			err := actual.Scan(test.input)

			assert.ErrorIs(t, err, domainerrors.ErrUnsupportedDatabaseEnumValue)
			assert.Equal(t, models.LinearScaling, actual.ToDomain())
		})
	}
}

func TestIT_DbScalingMode(t *testing.T) {
	tests := []struct {
		scaling  string
		expected models.ScalingMode
	}{
		{
			scaling:  "LINEAR",
			expected: models.LinearScaling,
		},
		{
			scaling:  "GEOMETRIC",
			expected: models.GeometricScaling,
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("reads scaling mode %v", test.scaling)
		t.Run(name, func(t *testing.T) {
			conn := newTestConnection(t)

			buildingId := insertTestBuilding(t, conn)

			_, err := conn.Exec(
				t.Context(),
				`INSERT INTO building_resource_metabolization_ship_speedup
						(building, scaling, base, coefficient)
						VALUES ($1, $2, $3, $4)`,
				buildingId,
				test.scaling,
				1.0,
				1.0,
			)
			require.NoError(t, err)

			tx, err := conn.BeginTx(t.Context())
			require.NoError(t, err)
			defer tx.Close(t.Context())

			row, err := db.QueryOneTx[struct{ Scaling DbScalingMode }](
				t.Context(),
				tx,
				`SELECT scaling
					 FROM building_resource_metabolization_ship_speedup
					 WHERE building = $1`,
				buildingId,
			)
			require.NoError(t, err)
			assert.Equal(t, test.expected, row.Scaling.ToDomain())
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

		assert.ErrorContains(t, err, "violates check constraint")
		assert.ErrorContains(t, err, "SQLSTATE 23514")
	})
}

func newTestConnection(t *testing.T) db.Connection {
	t.Helper()
	return sharedDbContainer.NewTestConnection(t)
}

func insertTestBuilding(
	t *testing.T,
	conn db.Connection,
) uuid.UUID {
	t.Helper()

	buildingId := uuid.New()
	name := fmt.Sprintf("scaling-mode-test-%v", buildingId.String())

	_, err := conn.Exec(
		t.Context(),
		`INSERT INTO building (id, name) VALUES ($1, $2)`,
		buildingId,
		name,
	)
	require.NoError(t, err, "Actual err: %v", err)

	return buildingId
}
