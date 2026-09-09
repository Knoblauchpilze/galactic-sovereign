package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_BuildingShipSpeedup_FactorAt(t *testing.T) {
	type testCase struct {
		name           string
		speedup        BuildingShipSpeedup
		level          int
		expectedFactor float64
	}

	testCases := []testCase{
		{
			name: "linear level 0",
			speedup: BuildingShipSpeedup{
				Scaling:  LinearScaling,
				Base:     0.47,
				Progress: 1.59,
			},
			level:          0,
			expectedFactor: 0.47,
		},
		{
			name: "linear level 2",
			speedup: BuildingShipSpeedup{
				Scaling:  LinearScaling,
				Base:     1.17,
				Progress: 0.23,
			},
			level:          2,
			expectedFactor: 1.63,
		},
		{
			name: "geometric level 0",
			speedup: BuildingShipSpeedup{
				Scaling:  GeometricScaling,
				Base:     0.2,
				Progress: 0.95,
			},
			level:          0,
			expectedFactor: 1.2,
		},
		{
			name: "geometric level 4",
			speedup: BuildingShipSpeedup{
				Scaling:  GeometricScaling,
				Base:     1.2,
				Progress: 2.85,
			},
			level:          4,
			expectedFactor: 67.17500625000001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.speedup.FactorAt(tc.level)
			assert.Equal(t, tc.expectedFactor, actual)
		})
	}
}
