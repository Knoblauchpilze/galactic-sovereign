package models

import (
	"testing"
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	someTime = time.Date(2026, 6, 8, 8, 22, 35, 0, time.UTC)
)

func TestUnit_Building_CreateBuildingAction(t *testing.T) {
	buildingId := uuid.New()
	sampleResource1 := uuid.New()
	sampleResource2 := uuid.New()
	p := Planet{
		Id: uuid.New(),
		// High enough values to not have to worry about costs
		Resources: []PlanetResource{
			{
				Resource: sampleResource1,
				Amount:   999999,
			},
			{
				Resource: sampleResource2,
				Amount:   999999,
			},
			{
				Resource: metalResourceId,
				Amount:   999999,
			},
			{
				Resource: crystalResourceId,
				Amount:   999999,
			},
		},
		Buildings: []PlanetBuilding{
			{
				Building: buildingId,
				Level:    4,
			},
		},
	}

	t.Run("correctly calculates action costs", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource: sampleResource1,
					Cost:     36,
					Progress: 1.5,
				},
				{
					Resource: sampleResource2,
					Cost:     78,
					Progress: 1.7,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p)
		require.NoError(t, err, "Actual err: %v", err)

		expected := BuildingAction{
			// The identifier is generated
			Id:           action.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,

			CreatedAt: action.CreatedAt,
			// Ignore the completion here, there are dedicated tests
			CompletedAt: action.CompletedAt,

			Version: 0,

			Costs: []BuildingActionCost{
				{
					Resource: sampleResource1,
					Amount:   182,
				},
				{
					Resource: sampleResource2,
					Amount:   651,
				},
			},
			Storages:    []BuildingActionResourceStorage{},
			Productions: []BuildingActionResourceProduction{},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("correctly calculates action resource productions", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs:     []BuildingCost{},
			Productions: []BuildingResourceProduction{
				{
					Resource: sampleResource1,
					Base:     74,
					Progress: 4.5,
				},
				{
					Resource: sampleResource2,
					Base:     98,
					Progress: 2.4,
				},
			},
			Storages: []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p)
		require.NoError(t, err, "Actual err: %v", err)

		expected := BuildingAction{
			Id:           action.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,

			CreatedAt: action.CreatedAt,
			// Ignore the completion here, there are dedicated tests
			CompletedAt: action.CompletedAt,

			Version: 0,

			Costs: []BuildingActionCost{},
			Productions: []BuildingActionResourceProduction{
				{
					Resource:   sampleResource1,
					Production: 682754,
				},
				{
					Resource:   sampleResource2,
					Production: 39016,
				},
			},
			Storages: []BuildingActionResourceStorage{},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("correctly calculates action resource storages", func(t *testing.T) {
		b := Building{
			Id:          buildingId,
			Name:        "test-building",
			CreatedAt:   someTime,
			Costs:       []BuildingCost{},
			Productions: []BuildingResourceProduction{},
			Storages: []BuildingResourceStorage{
				{
					Resource: sampleResource1,
					Base:     8904,
					Scale:    2,
					Progress: 2.2,
				},
				{
					Resource: sampleResource2,
					Base:     312,
					Scale:    1.2,
					Progress: 1.1,
				},
			},
		}

		action, err := b.CreateBuildingAction(p)
		require.NoError(t, err, "Actual err: %v", err)

		expected := BuildingAction{
			Id:           action.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,

			CreatedAt: action.CreatedAt,
			// Ignore the completion here, there are dedicated tests
			CompletedAt: action.CompletedAt,

			Version: 0,

			Costs:       []BuildingActionCost{},
			Productions: []BuildingActionResourceProduction{},
			Storages: []BuildingActionResourceStorage{
				{
					Resource: sampleResource1,
					Storage:  917112,
				},
				{
					Resource: sampleResource2,
					Storage:  312,
				},
			},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("correctly calculates completion time when no metal nor crystal is used", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource: sampleResource1,
					Cost:     36,
					Progress: 1.5,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p)
		require.NoError(t, err, "Actual err: %v", err)

		actual := action.CompletedAt.Sub(action.CreatedAt)
		assert.Equal(t, time.Duration(0), actual)
		expectedCosts := []BuildingActionCost{
			{
				Resource: sampleResource1,
				Amount:   182,
			},
		}
		assert.Equal(t, expectedCosts, action.Costs)
	})

	t.Run("correctly calculates completion time when no resource is used", func(t *testing.T) {
		b := Building{
			Id:          buildingId,
			Name:        "test-building",
			CreatedAt:   someTime,
			Costs:       []BuildingCost{},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p)
		require.NoError(t, err, "Actual err: %v", err)

		actual := action.CompletedAt.Sub(action.CreatedAt)
		assert.Equal(t, time.Duration(0), actual)
	})

	t.Run("correctly calculates completion time for metal usage", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource: metalResourceId,
					Cost:     36,
					Progress: 1.5,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p)
		require.NoError(t, err, "Actual err: %v", err)

		actual := action.CompletedAt.Sub(action.CreatedAt)
		completionTime := 262080 * time.Millisecond
		assert.Equal(t, completionTime, actual)
	})

	t.Run("correctly calculates completion time for crystal usage", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource: crystalResourceId,
					Cost:     79,
					Progress: 1.7,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p)
		require.NoError(t, err, "Actual err: %v", err)

		actual := action.CompletedAt.Sub(action.CreatedAt)
		completionTime := 948960 * time.Millisecond
		assert.Equal(t, completionTime, actual)
	})

	t.Run("correctly calculates completion time when metal and crystal are used", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource: metalResourceId,
					Cost:     36,
					Progress: 1.5,
				},
				{
					Resource: crystalResourceId,
					Cost:     79,
					Progress: 1.7,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p)
		require.NoError(t, err, "Actual err: %v", err)

		actual := action.CompletedAt.Sub(action.CreatedAt)
		completionTime := 1211040 * time.Millisecond
		assert.Equal(t, completionTime, actual)
	})

	t.Run("returns error when building does not exist on planet", func(t *testing.T) {
		b := Building{Id: uuid.New()}

		_, err := b.CreateBuildingAction(p)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingNotFound)
	})

	t.Run("returns error when resources are missing on the planet", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource: metalResourceId,
					Cost:     36,
					Progress: 1.5,
				},
				{
					Resource: crystalResourceId,
					Cost:     79,
					Progress: 1.7,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		poorPlanet := p
		poorPlanet.Resources = []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   189,
			},
			{
				Resource: crystalResourceId,
				// Needed value: 659
				Amount: 658,
			},
		}

		_, err := b.CreateBuildingAction(poorPlanet)

		assert.ErrorIs(t, domainerrors.ErrNotEnoughResources, err)
	})
}
