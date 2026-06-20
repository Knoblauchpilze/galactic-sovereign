package models

import (
	"slices"
	"testing"
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_Planet_AddBuildingAction(t *testing.T) {
	t.Run("returns error when planet already has an action", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		actionId := uuid.New()
		p.BuildingAction = &actionId

		b := generateTestBuilding(t)

		_, err := p.AddBuildingAction(b)

		assert.ErrorIs(t, domainerrors.ErrActionAlreadyInProgress, err)
	})

	t.Run("returns error when planet does not have enough resources", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)
		p.Resources = []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   189,
			},
			{
				Resource: crystalResourceId,
				// Needed value: 651
				Amount: 650,
			},
		}

		b := generateTestBuilding(t, withBuildingCost)

		_, err := p.AddBuildingAction(b)

		assert.ErrorIs(t, err, domainerrors.ErrNotEnoughResources, "Actual err: %v", err)
	})

	t.Run("returns error when building does not exist on planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)

		b := Building{Id: uuid.New()}

		_, err := p.AddBuildingAction(b)

		assert.ErrorIs(t, domainerrors.ErrBuildingNotFound, err)
	})

	t.Run("correctly assigns building action to planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		action, err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, &action.Id, p.BuildingAction)
	})

	t.Run("returns expected action", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		action, err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)

		completionTime := 1199520 * time.Millisecond
		expectedAction := BuildingAction{
			Id:           action.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,
			CreatedAt:    action.CreatedAt,
			CompletedAt:  action.CreatedAt.Add(completionTime),
			Version:      0,
			Costs: []BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   182,
				},
				{
					Resource: crystalResourceId,
					Amount:   651,
				},
			},
			Storages: []BuildingActionResourceStorage{
				{
					Resource: metalResourceId,
					Storage:  917112,
				},
				{
					Resource: crystalResourceId,
					Storage:  312,
				},
			},
			Productions: []BuildingActionResourceProduction{
				{
					Resource:   metalResourceId,
					Production: 682754,
				},
				{
					Resource:   crystalResourceId,
					Production: 39016,
				},
			},
		}
		assert.Equal(t, expectedAction, action)
	})

	t.Run("deducts action costs from the available planet resources", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		initialResources := slices.Clone(p.Resources)

		action, err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)

		expectedResources := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   initialResources[0].Amount - float64(action.Costs[0].Amount),
			},
			{
				Resource: crystalResourceId,
				Amount:   initialResources[1].Amount - float64(action.Costs[1].Amount),
			},
		}
		assert.Equal(t, expectedResources, p.Resources)
	})

	t.Run("bumps version and updated at field", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.UpdatedAt = someTime
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		initialVersion := p.Version

		action, err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, initialVersion+1, p.Version)
		assert.Equal(t, action.CreatedAt, p.UpdatedAt)
	})
}

func generateTestPlanet(
	t *testing.T,
	modifiers ...func(*testing.T, *Planet),
) Planet {
	t.Helper()

	p := Planet{
		Id: uuid.New(),
	}

	for _, modifier := range modifiers {
		modifier(t, &p)
	}

	return p
}

func withManyResources(t *testing.T, p *Planet) {
	t.Helper()

	// High enough values to not have to worry about costs
	p.Resources = []PlanetResource{
		{
			Resource: metalResourceId,
			Amount:   999999,
		},
		{
			Resource: crystalResourceId,
			Amount:   999999,
		},
	}
}

func withPlanetBuilding(t *testing.T, p *Planet) {
	p.Buildings = []PlanetBuilding{
		{
			Building: buildingId,
			Level:    4,
		},
	}
}
