package persistence

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var resourceId = uuid.MustParse("e0e56162-b462-4c05-9bda-828a373037a7")

func TestToPlanetResourceProduction(t *testing.T) {
	assert := assert.New(t)

	actionProduction := BuildingActionResourceProduction{
		Action:     actionId,
		Resource:   resourceId,
		Production: 54,
	}
	planetProduction := PlanetResourceProduction{
		Planet:     planetId,
		Building:   &buildingId,
		Resource:   resourceId,
		Production: 59,

		CreatedAt: time.Date(2024, 9, 14, 20, 40, 05, 651387249, time.UTC),
		UpdatedAt: time.Date(2024, 9, 14, 20, 40, 06, 651387249, time.UTC),

		Version: 1,
	}

	actual := ToPlanetResourceProduction(actionProduction, planetProduction)

	assert.Equal(planetId, actual.Planet)
	assert.Equal(&buildingId, actual.Building)
	assert.Equal(resourceId, actual.Resource)
	assert.Equal(actionProduction.Production, actual.Production)
	assert.Equal(planetProduction.CreatedAt, actual.CreatedAt)
	assert.Equal(planetProduction.UpdatedAt, actual.UpdatedAt)
	assert.Equal(planetProduction.Version, actual.Version)
}
