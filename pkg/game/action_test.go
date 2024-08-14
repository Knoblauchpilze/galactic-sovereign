package game

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultId = uuid.MustParse("60043fb4-d4bc-4bf0-95fd-dcdaf09a6acc")

func TestConsolidateBuildingAction_WhenNoBuilding_SetsDefault(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	buildings := []persistence.PlanetBuilding{}
	costs := []persistence.BuildingCost{}

	actual, err := ConsolidateBuildingAction(action, buildings, defaultResources, costs)

	assert.Nil(err)
	assert.Equal(0, actual.CurrentLevel)
	assert.Equal(1, actual.DesiredLevel)
}

func TestConsolidateBuildingAction_SetsCompletionTime(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	buildings := []persistence.PlanetBuilding{}
	costs := []persistence.BuildingCost{
		{
			Resource: defaultMetalId,
			Cost:     1250,
		},
		{
			Resource: defaultCrystalId,
			Cost:     3750,
		},
	}

	actual, err := ConsolidateBuildingAction(action, buildings, defaultResources, costs)

	assert.Nil(err)
	expectedCompletionTime := actual.CreatedAt.Add(2 * time.Hour)
	assert.Equal(expectedCompletionTime, actual.CompletedAt)
}

func TestConsolidateBuildingAction(t *testing.T) {
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
	costs := []persistence.BuildingCost{}

	actual, err := ConsolidateBuildingAction(action, buildings, defaultResources, costs)

	assert.Nil(err)
	assert.Equal(26, actual.CurrentLevel)
	assert.Equal(27, actual.DesiredLevel)
}

func TestValidateActionBuilding_NoSuchBuilding(t *testing.T) {
	assert := assert.New(t)

	action := persistence.BuildingAction{
		Building: defaultId,
	}
	buildings := []persistence.PlanetBuilding{}

	err := validateActionBuilding(action, buildings)

	assert.True(errors.IsErrorWithCode(err, NoSuchBuilding))
}

func TestValidateActionBuilding(t *testing.T) {
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

func TestValidateBuildingAction_BuildingUnknown(t *testing.T) {
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
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     4,
		},
	}

	err := ValidateBuildingAction(action, resources, buildings, costs)

	assert.True(errors.IsErrorWithCode(err, NoSuchBuilding))
}

func TestValidateBuildingAction_CostFails(t *testing.T) {
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
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     4,
		},
	}

	err := ValidateBuildingAction(action, resources, buildings, costs)

	assert.True(errors.IsErrorWithCode(err, NotEnoughResources))
}

func TestValidateBuildingAction(t *testing.T) {
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
	costs := []persistence.BuildingCost{
		{
			Resource: defaultId,
			Cost:     1,
		},
	}

	err := ValidateBuildingAction(action, resources, buildings, costs)

	assert.Nil(err)
}
