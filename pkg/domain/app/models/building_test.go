package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	someTime = time.Date(2026, time.June, 8, 8, 22, 35, 0, time.UTC)

	buildingId = uuid.New()
)

func TestUnit_Building_CreateBuildingAction(t *testing.T) {
	t.Run("correctly calculates action costs", func(t *testing.T) {
		b := generateTestBuilding(t, withBuildingCost)

		action := b.CreateBuildingAction(5, someTime)

		expected := BuildingAction{
			// The identifier is generated
			Id:           action.Id,
			Building:     b.Id,
			DesiredLevel: 5,

			CreatedAt: someTime,
			// Ignore the completion here, there are dedicated tests
			CompletedAt: action.CompletedAt,

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

		action := b.CreateBuildingAction(5, someTime)

		expected := BuildingAction{
			Id:           action.Id,
			Building:     b.Id,
			DesiredLevel: 5,

			CreatedAt: someTime,
			// Ignore the completion time here, there are dedicated tests
			CompletedAt: action.CompletedAt,

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

		action := b.CreateBuildingAction(5, someTime)

		expected := BuildingAction{
			Id:           action.Id,
			Building:     b.Id,
			DesiredLevel: 5,

			CreatedAt: someTime,
			// Ignore the completion time here, there are dedicated tests
			CompletedAt: action.CompletedAt,

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

	t.Run("correctly calculates completion time based on build time per unit", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource:              uuid.New(),
					Cost:                  36,
					Progress:              1.5,
					BuildTimeHoursPerUnit: 1,
				},
				{
					Resource:              uuid.New(),
					Cost:                  15,
					Progress:              1.2,
					BuildTimeHoursPerUnit: 36,
				},
				{
					Resource:              uuid.New(),
					Cost:                  100,
					Progress:              1.8,
					BuildTimeHoursPerUnit: 0.04,
				},
				{
					Resource:              uuid.New(),
					Cost:                  150,
					Progress:              1.01,
					BuildTimeHoursPerUnit: 0,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action := b.CreateBuildingAction(5, someTime)

		expectedCosts := []BuildingActionCost{
			{
				Resource: b.Costs[0].Resource,
				Amount:   182,
			},
			{
				Resource: b.Costs[1].Resource,
				Amount:   31,
			},
			{
				Resource: b.Costs[2].Resource,
				Amount:   1049,
			},
			{
				Resource: b.Costs[3].Resource,
				Amount:   156,
			},
		}
		assert.Equal(t, expectedCosts, action.Costs)

		completionTime := 4823856 * time.Second
		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime.Add(completionTime), action.CompletedAt)
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

		action := b.CreateBuildingAction(5, someTime)

		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime, action.CompletedAt)
	})

	t.Run("correctly calculates completion time when single resource is used", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource:              uuid.New(),
					Cost:                  36,
					Progress:              1.5,
					BuildTimeHoursPerUnit: 0.0004,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action := b.CreateBuildingAction(5, someTime)

		completionTime := 262080 * time.Millisecond
		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime.Add(completionTime), action.CompletedAt)
	})

	t.Run("correctly calculates completion time when resource has no build time", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource:              crystalResourceId,
					Cost:                  79,
					Progress:              1.7,
					BuildTimeHoursPerUnit: 0.0,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action := b.CreateBuildingAction(5, someTime)

		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime, action.CompletedAt)
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
			Resource:              metalResourceId,
			Cost:                  36,
			Progress:              1.5,
			BuildTimeHoursPerUnit: 0.0004,
		},
		{
			Resource:              crystalResourceId,
			Cost:                  78,
			Progress:              1.7,
			BuildTimeHoursPerUnit: 0.0004,
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
