package game

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultId = uuid.MustParse("60043fb4-d4bc-4bf0-95fd-dcdaf09a6acc")

func TestUnit_DetermineBuildingActionCost(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Id:           uuid.MustParse("7f548f48-2bac-46f0-b655-56487472b5db"),
		DesiredLevel: 2,
	}
	baseCosts := []persistence.BuildingCost{
		{
			Building: defaultId,
			Resource: defaultMetalId,
			Cost:     3,
			Progress: 1.2,
		},
		{
			Building: defaultId,
			Resource: defaultCrystalId,
			Cost:     6,
			Progress: 1.2,
		},
	}

	costs := DetermineBuildingActionCost(action, baseCosts)

	assert.Equal(2, len(costs))
	expectedCost := persistence.BuildingActionCost{
		Action:   action.Id,
		Resource: defaultMetalId,
		Amount:   3,
	}
	assert.Equal(expectedCost, costs[0])
	expectedCost = persistence.BuildingActionCost{
		Action:   action.Id,
		Resource: defaultCrystalId,
		Amount:   7,
	}
	assert.Equal(expectedCost, costs[1])
}

func TestUnit_DetermineBuildingActionResourceProduction(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Id:           uuid.MustParse("7f548f48-2bac-46f0-b655-56487472b5db"),
		DesiredLevel: 10,
	}
	baseProductions := []persistence.BuildingResourceProduction{
		{
			Building: defaultId,
			Resource: defaultMetalId,
			Base:     21,
			Progress: 1.2,
		},
		{
			Building: defaultId,
			Resource: defaultCrystalId,
			Base:     27,
			Progress: 1.3,
		},
	}

	productions := DetermineBuildingActionResourceProduction(action, baseProductions)

	assert.Equal(2, len(productions))
	expectedResourceProduction := persistence.BuildingActionResourceProduction{
		Action:     action.Id,
		Resource:   defaultMetalId,
		Production: 975,
	}
	assert.Equal(expectedResourceProduction, productions[0])
	expectedResourceProduction = persistence.BuildingActionResourceProduction{
		Action:     action.Id,
		Resource:   defaultCrystalId,
		Production: 2576,
	}
	assert.Equal(expectedResourceProduction, productions[1])
}

func TestUnit_ConsolidateBuildingActionLevel_WhenNoBuilding_SetsDefault(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	buildings := []persistence.PlanetBuilding{}

	actual := ConsolidateBuildingActionLevel(action, buildings)

	assert.Equal(0, actual.CurrentLevel)
	assert.Equal(1, actual.DesiredLevel)
}

func TestUnit_ConsolidateBuildingActionLevel_WhenBuildingExists_SetsCorrectLevel(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
			Level:    26,
		},
	}

	actual := ConsolidateBuildingActionLevel(action, buildings)

	assert.Equal(26, actual.CurrentLevel)
	assert.Equal(27, actual.DesiredLevel)
}

func TestUnit_ConsolidateBuildingActionCompletionTime_WhenResourceNotFound_ExpectError(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	resources := []persistence.Resource{
		{
			Id:   defaultMetalId,
			Name: "metal",
		},
	}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultMetalId,
			Amount:   1250,
		},
		{
			Resource: defaultCrystalId,
			Amount:   3750,
		},
	}

	_, err := ConsolidateBuildingActionCompletionTime(action, resources, costs)

	assert.True(errors.IsErrorWithCode(err, NoSuchResource))
}

func TestUnit_ConsolidateBuildingActionCompletionTime(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultMetalId,
			Amount:   1250,
		},
		{
			Resource: defaultCrystalId,
			Amount:   3750,
		},
	}

	actual, err := ConsolidateBuildingActionCompletionTime(action, defaultResources, costs)

	assert.Nil(err)
	expectedCompletionTime := actual.CreatedAt.Add(2 * time.Hour)
	assert.Equal(expectedCompletionTime, actual.CompletedAt)
}

func TestUnit_ValidateActionBuilding_NoSuchBuilding(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	buildings := []persistence.PlanetBuilding{}

	err := validateActionBuilding(action, buildings)

	assert.True(errors.IsErrorWithCode(err, NoSuchBuilding))
}

func TestUnit_ValidateActionBuilding(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
		},
	}

	err := validateActionBuilding(action, buildings)

	assert.Nil(err)
}

func TestUnit_ValidateActionCost_NoSuchResource(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.PlanetResource{}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultId,
			Amount:   1,
		},
	}

	err := validateActionCost(resources, costs)

	assert.True(errors.IsErrorWithCode(err, NotEnoughResources))
}

func TestUnit_ValidateActionCost_TooLittleResource(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   1,
		},
	}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}

	err := validateActionCost(resources, costs)

	assert.True(errors.IsErrorWithCode(err, NotEnoughResources))
}

func TestUnit_ValidateActionCost_ExactlyEnoughResource(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}

	err := validateActionCost(resources, costs)

	assert.Nil(err)
}

func TestUnit_ValidateActionCost_MoreThanEnoughResource(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2.5,
		},
	}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}

	err := validateActionCost(resources, costs)

	assert.Nil(err)
}

func TestUnit_ValidateBuildingAction_BuildingUnknown(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	buildings := []persistence.PlanetBuilding{}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultId,
			Amount:   4,
		},
	}

	err := ValidateBuildingAction(action, resources, buildings, costs)

	assert.True(errors.IsErrorWithCode(err, NoSuchBuilding))
}

func TestUnit_ValidateBuildingAction_CostFails(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
		},
	}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultId,
			Amount:   4,
		},
	}

	err := ValidateBuildingAction(action, resources, buildings, costs)

	assert.True(errors.IsErrorWithCode(err, NotEnoughResources))
}

func TestUnit_ValidateBuildingAction(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
		},
	}
	costs := []persistence.BuildingActionCost{
		{
			Resource: defaultId,
			Amount:   1,
		},
	}

	err := ValidateBuildingAction(action, resources, buildings, costs)

	assert.Nil(err)
}
