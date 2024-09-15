package game

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultPlanetResource = persistence.PlanetResource{
	Planet:   uuid.MustParse("ab6806b1-722b-4438-afa4-2b23a389d773"),
	Resource: uuid.MustParse("c55df64d-5df1-4acf-bcca-f0a5e5749d87"),
	Amount:   60.0,

	UpdatedAt: time.Date(2024, 9, 13, 14, 51, 55, 651387248, time.UTC),

	Version: 10,
}
var resourceProductionPerHour = 5.0

var thresholdForResourceEquality = 1e-6

func generatePlanetResource() persistence.PlanetResource {
	resource := defaultPlanetResource
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()

	return resource
}

func TestUpdatePlanetResourceAmountToTime_whenTimeInThePast_expectNoUpdate(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	inThePast := resource.UpdatedAt.Add(-1 * time.Hour)
	updated := UpdatePlanetResourceAmountToTime(resource, resourceProductionPerHour, inThePast)

	assert.Equal(resource.Amount, updated.Amount)
}

func TestUpdatePlanetResourceAmountToTime_updatesAmount(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	type testCase struct {
		duration       time.Duration
		expectedAmount float64
	}

	testCases := map[string]testCase{
		"1s": {
			duration:       1 * time.Second,
			expectedAmount: 60.001389,
		},
		"41s": {
			duration:       41 * time.Second,
			expectedAmount: 60.056944,
		},
		"1m": {
			duration:       1 * time.Minute,
			expectedAmount: 60.083333,
		},
		"30m": {
			duration:       30 * time.Minute,
			expectedAmount: 62.5,
		},
		"1h": {
			duration:       1 * time.Hour,
			expectedAmount: 65,
		},
		"4h3m2s2ms": {
			duration:       4*time.Hour + 3*time.Minute + 2*time.Second + 1*time.Millisecond,
			expectedAmount: 80.252779,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			oneHourFromNow := resource.UpdatedAt.Add(testCase.duration)
			updated := UpdatePlanetResourceAmountToTime(resource, resourceProductionPerHour, oneHourFromNow)

			assert.InDelta(testCase.expectedAmount, updated.Amount, thresholdForResourceEquality)
		})
	}
}

func TestUpdatePlanetResourceAmountToTime_updatesUpdatedAt(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	oneHourFromNow := resource.UpdatedAt.Add(1 * time.Hour)
	updated := UpdatePlanetResourceAmountToTime(resource, resourceProductionPerHour, oneHourFromNow)

	assert.Equal(oneHourFromNow, updated.UpdatedAt)
}
