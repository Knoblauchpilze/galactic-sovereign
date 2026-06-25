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
		p.BuildingAction = &BuildingAction{Id: actionId}

		b := generateTestBuilding(t)

		err := p.AddBuildingAction(b)

		assert.ErrorIs(t, domainerrors.ErrActionAlreadyInProgress, err)
		require.NotNil(t, p.BuildingAction)
		assert.Equal(t, actionId, p.BuildingAction.Id)
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

		err := p.AddBuildingAction(b)

		assert.ErrorIs(t, err, domainerrors.ErrNotEnoughResources, "Actual err: %v", err)
		assert.Nil(t, p.BuildingAction)
	})

	t.Run("returns error when building does not exist on planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)

		b := Building{Id: uuid.New()}

		err := p.AddBuildingAction(b)

		assert.ErrorIs(t, domainerrors.ErrBuildingNotFound, err)
		assert.Nil(t, p.BuildingAction)
	})

	t.Run("assigns building action to planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, p.BuildingAction)

		completionTime := 1199520 * time.Millisecond
		expectedAction := &BuildingAction{
			Id:           p.BuildingAction.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,
			CreatedAt:    p.BuildingAction.CreatedAt,
			CompletedAt:  p.BuildingAction.CreatedAt.Add(completionTime),
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
		assert.Equal(t, expectedAction, p.BuildingAction)
	})

	t.Run("deducts action costs from the available planet resources", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		initialResources := slices.Clone(p.Resources)

		err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, p.BuildingAction)

		expectedMetalAmount := initialResources[0].Amount - float64(p.BuildingAction.Costs[0].Amount)
		expectedCrystalAmount := initialResources[1].Amount - float64(p.BuildingAction.Costs[1].Amount)
		expectedResources := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   expectedMetalAmount,
			},
			{
				Resource: crystalResourceId,
				Amount:   expectedCrystalAmount,
			},
		}
		assert.Equal(t, expectedResources, p.Resources)
	})

	t.Run("bumps version and updated at field", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.UpdatedAt = someTime
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		initialVersion := p.Version

		err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, p.BuildingAction)

		assert.Equal(t, initialVersion+1, p.Version)
		assert.Equal(t, p.BuildingAction.CreatedAt, p.UpdatedAt)
	})
}

func TestUnit_Planet_CancelBuildingAction(t *testing.T) {
	t.Run("returns error when planet does not have an action", func(t *testing.T) {
	})

	t.Run("resets building action in planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.BuildingAction = &BuildingAction{Id: uuid.New()}

		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, p.BuildingAction)
	})

	t.Run("adds back action costs to the available planet resources", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)
		p.Resources = []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   1000,
			},
			{
				Resource: crystalResourceId,
				Amount:   2000,
			},
		}

		p.BuildingAction = &BuildingAction{
			Id: uuid.New(),
			Costs: []BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   36,
				},
				{
					Resource: crystalResourceId,
					Amount:   178,
				},
			},
		}

		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   1036,
			},
			{
				Resource: crystalResourceId,
				Amount:   2178,
			},
		}
		assert.Equal(t, expected, p.Resources)
	})

	t.Run("adds back action costs when resources are not present on the planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)
		p.Resources = []PlanetResource{}
		p.BuildingAction = &BuildingAction{
			Id: uuid.New(),
			Costs: []BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   36,
				},
				{
					Resource: crystalResourceId,
					Amount:   178,
				},
			},
		}

		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   36,
			},
			{
				Resource: crystalResourceId,
				Amount:   178,
			},
		}
		assert.Equal(t, expected, p.Resources)
	})

	t.Run("bumps version and updated at field", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.BuildingAction = &BuildingAction{Id: uuid.New()}

		initialVersion := p.Version

		beforeCall := time.Now()
		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, initialVersion+1, p.Version)
		assert.True(t, p.UpdatedAt.After(beforeCall))
	})
}

func TestUnit_Planet_UpdateToTime(t *testing.T) {
	t.Run("does not update resource when no production is defined", func(t *testing.T) {})

	t.Run("does not change resource value when already over the storage capacity", func(t *testing.T) {})

	t.Run("does not change resource when update time is before already computed time", func(t *testing.T) {})

	t.Run("caps resource at storage capacity", func(t *testing.T) {})

	t.Run("adds full production value when storage is sufficient", func(t *testing.T) {})

	t.Run("adds all production for a resource", func(t *testing.T) {})

	t.Run("adds a resource when it is produced but not yet stored on the planet", func(t *testing.T) {})

	t.Run("keeps resource at 0 when no storage is defined for it", func(t *testing.T) {})

	t.Run("bumps updated at field", func(t *testing.T) {})
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
