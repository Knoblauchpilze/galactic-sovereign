package game

import (
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultId = uuid.MustParse("60043fb4-d4bc-4bf0-95fd-dcdaf09a6acc")

func TestValidateActionLevel_DesiredLevelIsTheSameAsCurrentLevel(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		CurrentLevel: 0,
		DesiredLevel: 0,
	}
	err := validateActionLevel(action)
	assert.True(errors.IsErrorWithCode(err, InvalidActionData))
}

func TestValidateActionLevel_DesiredLevelIsSmallerAsCurrentLevel(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		CurrentLevel: 1,
		DesiredLevel: 0,
	}
	err := validateActionLevel(action)
	assert.True(errors.IsErrorWithCode(err, InvalidActionData))
}

func TestValidateActionLevel_DesiredLevelIsTooBigComparedToCurrentLevel(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		CurrentLevel: 1,
		DesiredLevel: 3,
	}
	err := validateActionLevel(action)
	assert.True(errors.IsErrorWithCode(err, InvalidActionData))
}

func TestValidateActionLevel(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		CurrentLevel: 1,
		DesiredLevel: 2,
	}
	err := validateActionLevel(action)
	assert.Nil(err)
}

func TestValidateActionBuilding_NoSuchBuilding(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	buildings := []persistence.PlanetBuilding{}
	err := validateActionBuilding(action, buildings)
	assert.True(errors.IsErrorWithCode(err, InvalidBuildingLevel))
}

func TestValidateActionBuilding_BuildingExistsButWithADifferentLevel(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building:     defaultId,
		CurrentLevel: 1,
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
			Level:    2,
		},
	}
	err := validateActionBuilding(action, buildings)
	assert.True(errors.IsErrorWithCode(err, InvalidBuildingLevel))
}

func TestValidateActionBuilding(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building:     defaultId,
		CurrentLevel: 1,
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
			Level:    1,
		},
	}
	err := validateActionBuilding(action, buildings)
	assert.Nil(err)
}

func TestValidateActionCost_NoSuchResource(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.PlanetResource{}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     1,
		},
	}
	err := validateActionCost(resources, costs)
	assert.True(errors.IsErrorWithCode(err, NotEnoughResources))
}

func TestValidateActionCost_TooLittleResource(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   1,
		},
	}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     2,
		},
	}
	err := validateActionCost(resources, costs)
	assert.True(errors.IsErrorWithCode(err, NotEnoughResources))
}

func TestValidateActionCost_ExactlyEnoughResource(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     2,
		},
	}
	err := validateActionCost(resources, costs)
	assert.Nil(err)
}

func TestValidateActionCost_MoreThanEnoughResource(t *testing.T) {
	assert := assert.New(t)

	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2.5,
		},
	}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     2,
		},
	}
	err := validateActionCost(resources, costs)
	assert.Nil(err)
}

func TestValidateBuildingAction_LevelFails(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building:     defaultId,
		CurrentLevel: 1,
		DesiredLevel: 3,
	}
	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     1,
		},
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
			Level:    1,
		},
	}
	err := ValidateBuildingAction(action, resources, costs, buildings)
	assert.True(errors.IsErrorWithCode(err, InvalidActionData))
}

func TestValidateBuildingAction_BuildingFails(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building:     defaultId,
		CurrentLevel: 1,
		DesiredLevel: 2,
	}
	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     1,
		},
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
			Level:    3,
		},
	}
	err := ValidateBuildingAction(action, resources, costs, buildings)
	assert.True(errors.IsErrorWithCode(err, InvalidBuildingLevel))
}

func TestValidateBuildingAction_CostFails(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building:     defaultId,
		CurrentLevel: 1,
		DesiredLevel: 2,
	}
	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     4,
		},
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
			Level:    1,
		},
	}
	err := ValidateBuildingAction(action, resources, costs, buildings)
	assert.True(errors.IsErrorWithCode(err, NotEnoughResources))
}

func TestValidateBuildingAction(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building:     defaultId,
		CurrentLevel: 1,
		DesiredLevel: 2,
	}
	resources := []persistence.PlanetResource{
		{
			Resource: defaultId,
			Amount:   2,
		},
	}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     1,
		},
	}
	buildings := []persistence.PlanetBuilding{
		{
			Building: defaultId,
			Level:    1,
		},
	}
	err := ValidateBuildingAction(action, resources, costs, buildings)
	assert.Nil(err)
}
