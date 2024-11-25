package persistence

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var actionId = uuid.MustParse("f9fe2719-fded-4254-861e-6bd672a8662a")
var planetId = uuid.MustParse("328c7dd2-7004-4a90-91c7-162ca9adf628")
var buildingId = uuid.MustParse("4afa0ce8-3da3-4112-b10f-935919e962c3")

func TestUnit_ToBuildingAction(t *testing.T) {
	assert := assert.New(t)

	action := BuildingAction{
		Id:           actionId,
		Planet:       planetId,
		Building:     buildingId,
		CurrentLevel: 6,
		DesiredLevel: 7,
		CreatedAt:    time.Date(2024, 8, 17, 14, 07, 18, 651387246, time.UTC),
		CompletedAt:  time.Date(2024, 8, 17, 14, 07, 19, 651387246, time.UTC),
	}
	building := PlanetBuilding{
		Planet:    planetId,
		Building:  buildingId,
		Level:     6,
		CreatedAt: time.Date(2024, 8, 17, 14, 07, 10, 651387246, time.UTC),
		UpdatedAt: time.Date(2024, 8, 17, 14, 07, 15, 651387246, time.UTC),
	}

	actual := ToPlanetBuilding(action, building)

	assert.Equal(planetId, actual.Planet)
	assert.Equal(buildingId, actual.Building)
	assert.Equal(action.DesiredLevel, actual.Level)
	assert.Equal(actual.CreatedAt, building.CreatedAt)
	assert.Equal(action.CompletedAt, actual.UpdatedAt)
	assert.Equal(building.Version, actual.Version)
}
