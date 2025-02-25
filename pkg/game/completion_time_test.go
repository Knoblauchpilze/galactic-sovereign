package game

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultMetalId = uuid.MustParse("24c2a21a-3de8-42dd-bee4-8652e8368a5c")
var defaultCrystalId = uuid.MustParse("4e8a8ee5-668e-42c4-a4a5-938e5a68741c")

var defaultResources = []persistence.Resource{
	{
		Id:   defaultMetalId,
		Name: "metal",
	},
	{
		Id:   defaultCrystalId,
		Name: "crystal",
	},
}

func TestUnit_BuildingCompletionTimeFromCost_whenMetalNotFound_expectError(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.Resource{}
	costs := []persistence.BuildingActionCost{}

	_, err := buildingCompletionTimeFromCost(resources, costs)

	assert.True(errors.IsErrorWithCode(err, NoSuchResource))
}

func TestUnit_BuildingCompletionTimeFromCost_whenCrystalNotFound_expectError(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.Resource{
		{
			Id:   defaultMetalId,
			Name: "metal",
		},
	}
	costs := []persistence.BuildingActionCost{}

	_, err := buildingCompletionTimeFromCost(resources, costs)

	assert.True(errors.IsErrorWithCode(err, NoSuchResource))
}

func TestUnit_BuildingCompletionTimeFromCost_onlyMetalCost(t *testing.T) {
	assert := assert.New(t)

	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultMetalId,
			Amount:   1250,
		},
	}

	duration, err := buildingCompletionTimeFromCost(defaultResources, costs)

	assert.Nil(err)
	assert.Equal(30*time.Minute, duration)
}

func TestUnit_BuildingCompletionTimeFromCost_onlyCrystalCost(t *testing.T) {
	assert := assert.New(t)

	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultCrystalId,
			Amount:   3000,
		},
	}

	duration, err := buildingCompletionTimeFromCost(defaultResources, costs)

	assert.Nil(err)
	expectedDuration, err2 := time.ParseDuration("1h12m")
	assert.Nil(err2)
	assert.Equal(expectedDuration, duration)
}

func TestUnit_BuildingCompletionTimeFromCost_metalAndCrystal(t *testing.T) {
	assert := assert.New(t)

	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultMetalId,
			Amount:   5,
		},
		{
			Resource: defaultCrystalId,
			Amount:   5,
		},
	}

	duration, err := buildingCompletionTimeFromCost(defaultResources, costs)

	assert.Nil(err)
	expectedDuration, err2 := time.ParseDuration("14.4s")
	assert.Nil(err2)
	assert.Equal(expectedDuration, duration)
}
