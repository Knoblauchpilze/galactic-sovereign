package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultResourceId = uuid.MustParse("f568a280-92ec-4752-b14b-5e71070cce3e")
var otherResourceid = uuid.MustParse("966fe5c0-8309-4f62-853c-7d60ae5680f0")

var defaultPlanetResources = []persistence.PlanetResource{
	{
		Resource: defaultResourceId,
		Amount:   2.5,
	},
}
var defaultActionCosts = []persistence.BuildingActionCost{
	{
		Resource: defaultResourceId,
		Amount:   2,
	},
}

func TestFindResourceForCost_NoSuchResource(t *testing.T) {
	assert := assert.New(t)

	cost := persistence.BuildingActionCost{
		Resource: otherResourceid,
		Amount:   2,
	}

	_, err := findResourceForCost(defaultPlanetResources, cost)

	assert.True(errors.IsErrorWithCode(err, noSuchResource))
}

func TestFindResourceForCost_ResourceExist(t *testing.T) {
	assert := assert.New(t)

	cost := persistence.BuildingActionCost{
		Resource: defaultResourceId,
		Amount:   2,
	}

	actual, err := findResourceForCost(defaultPlanetResources, cost)

	assert.Nil(err)
	expected := defaultPlanetResources[0]
	assert.Equal(expected, actual)
}

func TestUpdatePlanetResourceWithCosts_ResourceNotFound(t *testing.T) {
	assert := assert.New(t)

	costs := []persistence.BuildingActionCost{
		{
			Resource: otherResourceid,
			Amount:   2,
		},
	}

	err := updatePlanetResourceWithCosts(context.Background(), nil, &mockPlanetResourceRepository{}, defaultPlanetResources, costs, addResource)

	assert.True(errors.IsErrorWithCode(err, FailedToCreateAction))
}

func TestUpdatePlanetResourceWithCosts_UpdateResourceInDb(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceRepository{}

	err := updatePlanetResourceWithCosts(context.Background(), nil, m, defaultPlanetResources, defaultActionCosts, addResource)

	assert.Nil(err)
	assert.Equal(1, m.updateCalled)
}

func TestUpdatePlanetResourceWithCosts_AddResource(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceRepository{}

	err := updatePlanetResourceWithCosts(context.Background(), nil, m, defaultPlanetResources, defaultActionCosts, addResource)

	assert.Nil(err)
	assert.Equal(1, len(m.updatedPlanetResources))
	assert.Equal(defaultPlanetResources[0].Resource, m.updatedPlanetResources[0].Resource)
	expectedAmount := defaultPlanetResources[0].Amount + float64(defaultActionCosts[0].Amount)
	assert.Equal(expectedAmount, m.updatedPlanetResources[0].Amount)
}

func TestUpdatePlanetResourceWithCosts_SubtractResource(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceRepository{}

	err := updatePlanetResourceWithCosts(context.Background(), nil, m, defaultPlanetResources, defaultActionCosts, subtractResource)

	assert.Nil(err)
	assert.Equal(1, len(m.updatedPlanetResources))
	assert.Equal(defaultPlanetResources[0].Resource, m.updatedPlanetResources[0].Resource)
	expectedAmount := defaultPlanetResources[0].Amount - float64(defaultActionCosts[0].Amount)
	assert.Equal(expectedAmount, m.updatedPlanetResources[0].Amount)
}

func TestUpdatePlanetResourceWithCosts_Update_Fails(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceRepository{
		updateErr: errDefault,
	}

	err := updatePlanetResourceWithCosts(context.Background(), nil, m, defaultPlanetResources, defaultActionCosts, addResource)

	assert.Equal(errDefault, err)
	assert.Equal(1, m.updateCalled)
}

func TestUpdatePlanetResourceWithCosts_Update_OptimisticLockException(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceRepository{
		updateErr: errors.NewCode(db.OptimisticLockException),
	}

	err := updatePlanetResourceWithCosts(context.Background(), nil, m, defaultPlanetResources, defaultActionCosts, addResource)

	assert.True(errors.IsErrorWithCode(err, ConflictingStateForAction))
	assert.Equal(1, m.updateCalled)
}
