package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnit_MergeWithPlanetResourceStorage(t *testing.T) {
	assert := assert.New(t)

	actionStorage := BuildingActionResourceStorage{
		Action:   actionId,
		Resource: resourceId,
		Storage:  9845,
	}
	planetStorage := PlanetResourceStorage{
		Planet:   planetId,
		Resource: resourceId,
		Storage:  4521,

		CreatedAt: time.Date(2025, 2, 12, 11, 43, 52, 651387249, time.UTC),
		UpdatedAt: time.Date(2025, 2, 12, 11, 43, 53, 651387249, time.UTC),

		Version: 1,
	}

	actual := MergeWithPlanetResourceStorage(actionStorage, planetStorage)

	assert.Equal(planetId, actual.Planet)
	assert.Equal(resourceId, actual.Resource)
	assert.Equal(actionStorage.Storage, actual.Storage)
	assert.Equal(planetStorage.CreatedAt, actual.CreatedAt)
	assert.Equal(planetStorage.UpdatedAt, actual.UpdatedAt)
	assert.Equal(planetStorage.Version, actual.Version)
}

func TestUnit_ToPlanetResourceStorage(t *testing.T) {
	assert := assert.New(t)

	actionStorage := BuildingActionResourceStorage{
		Action:   actionId,
		Resource: resourceId,
		Storage:  5689,
	}
	action := BuildingAction{
		Id:           actionId,
		Planet:       planetId,
		Building:     buildingId,
		CurrentLevel: 1,
		DesiredLevel: 2,

		CreatedAt:   time.Date(2025, 2, 12, 11, 45, 06, 651387251, time.UTC),
		CompletedAt: time.Date(2025, 2, 12, 11, 45, 07, 651387251, time.UTC),
	}

	actual := ToPlanetResourceStorage(actionStorage, action)

	assert.Equal(action.Planet, actual.Planet)
	assert.Equal(actionStorage.Resource, actual.Resource)
	assert.Equal(actionStorage.Storage, actual.Storage)
	assert.Equal(action.CompletedAt, actual.CreatedAt)
	assert.Equal(action.CompletedAt, actual.UpdatedAt)
	assert.Equal(0, actual.Version)
}
