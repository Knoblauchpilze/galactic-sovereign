package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	someTime = time.Date(2026, 6, 8, 8, 22, 35, 0, time.UTC)

	buildingId = uuid.New()
)

func TestUnit_Building_CreateBuildingAction(t *testing.T) {
	planetId := uuid.New()

	t.Run("correctly calculates action costs", func(t *testing.T) {
		b := generateTestBuilding(t, withBuildingCost)

		action := b.CreateBuildingAction(planetId, 5)

		expected := BuildingAction{
			// The identifier is generated
			Id:           action.Id,
			Planet:       planetId,
			Building:     b.Id,
			CurrentLevel: 4,
			DesiredLevel: 5,

			CreatedAt: action.CreatedAt,
			// Ignore the completion here, there are dedicated tests
			CompletedAt: action.CompletedAt,

			Version: 0,

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
			Storages:    []BuildingActionResourceStorage{},
			Productions: []BuildingActionResourceProduction{},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("correctly calculates action resource productions", func(t *testing.T) {
		b := generateTestBuilding(t, withBuildingProduction)

		action := b.CreateBuildingAction(planetId, 5)

		expected := BuildingAction{
			Id:           action.Id,
			Planet:       planetId,
			Building:     b.Id,
			CurrentLevel: 4,
			DesiredLevel: 5,

			CreatedAt: action.CreatedAt,
			// Ignore the completion here, there are dedicated tests
			CompletedAt: action.CompletedAt,

			Version: 0,

			Costs: []BuildingActionCost{},
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
			Storages: []BuildingActionResourceStorage{},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("correctly calculates action resource storages", func(t *testing.T) {
		b := generateTestBuilding(t, withBuildingStorage)

		action := b.CreateBuildingAction(planetId, 5)

		expected := BuildingAction{
			Id:           action.Id,
			Planet:       planetId,
			Building:     b.Id,
			CurrentLevel: 4,
			DesiredLevel: 5,

			CreatedAt: action.CreatedAt,
			// Ignore the completion here, there are dedicated tests
			CompletedAt: action.CompletedAt,

			Version: 0,

			Costs:       []BuildingActionCost{},
			Productions: []BuildingActionResourceProduction{},
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
					Resource: uuid.New(),
					Cost:     36,
					Progress: 1.5,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action := b.CreateBuildingAction(planetId, 5)

		actual := action.CompletedAt.Sub(action.CreatedAt)
		assert.Equal(t, time.Duration(0), actual)
		expectedCosts := []BuildingActionCost{
			{
				Resource: b.Costs[0].Resource,
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

		action := b.CreateBuildingAction(planetId, 5)

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

		action := b.CreateBuildingAction(planetId, 5)

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

		action := b.CreateBuildingAction(planetId, 5)

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

		action := b.CreateBuildingAction(planetId, 5)

		actual := action.CompletedAt.Sub(action.CreatedAt)
		completionTime := 1211040 * time.Millisecond
		assert.Equal(t, completionTime, actual)
	})
}

func generateTestBuilding(
	t *testing.T,
	modifiers ...func(*testing.T, *Building),
) Building {
	t.Helper()

	b := Building{
		Id:          buildingId,
		Name:        "test-building",
		CreatedAt:   someTime,
		Costs:       []BuildingCost{},
		Productions: []BuildingResourceProduction{},
		Storages:    []BuildingResourceStorage{},
	}

	for _, modifier := range modifiers {
		modifier(t, &b)
	}

	return b
}

func withBuildingCost(t *testing.T, b *Building) {
	t.Helper()

	b.Costs = []BuildingCost{
		{
			Resource: metalResourceId,
			Cost:     36,
			Progress: 1.5,
		},
		{
			Resource: crystalResourceId,
			Cost:     78,
			Progress: 1.7,
		},
	}
}

func withBuildingStorage(t *testing.T, b *Building) {
	t.Helper()

	b.Storages = []BuildingResourceStorage{
		{
			Resource: metalResourceId,
			Base:     8904,
			Scale:    2,
			Progress: 2.2,
		},
		{
			Resource: crystalResourceId,
			Base:     312,
			Scale:    1.2,
			Progress: 1.1,
		},
	}
}

func withBuildingProduction(t *testing.T, b *Building) {
	t.Helper()

	b.Productions = []BuildingResourceProduction{
		{
			Resource: metalResourceId,
			Base:     74,
			Progress: 4.5,
		},
		{
			Resource: crystalResourceId,
			Base:     98,
			Progress: 2.4,
		},
	}
}
