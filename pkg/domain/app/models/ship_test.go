package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	shipId = uuid.New()
)

func TestUnit_Ship_CreateShipAction(t *testing.T) {
	t.Run("correctly calculates action costs", func(t *testing.T) {
		s := generateTestShip(t, withShipCost)

		action := s.CreateShipAction(5, someTime)

		expected := ShipAction{
			// The identifier is generated
			Id:    action.Id,
			Ship:  s.Id,
			Count: 5,

			CreatedAt: someTime,
			// Ignore the completion here, there are dedicated tests
			NextCompletionAt: action.NextCompletionAt,
			CompletedAt:      action.CompletedAt,

			Costs: []ShipActionCost{
				{
					Resource: metalResourceId,
					Amount:   180,
				},
				{
					Resource: crystalResourceId,
					Amount:   390,
				},
			},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("correctly calculates costs based on build time per unit", func(t *testing.T) {
		s := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              uuid.New(),
					Cost:                  36,
					BuildTimeHoursPerUnit: 1,
				},
				{
					Resource:              uuid.New(),
					Cost:                  15,
					BuildTimeHoursPerUnit: 36,
				},
				{
					Resource:              uuid.New(),
					Cost:                  100,
					BuildTimeHoursPerUnit: 0.04,
				},
				{
					Resource:              uuid.New(),
					Cost:                  150,
					BuildTimeHoursPerUnit: 0,
				},
			},
		}

		action := s.CreateShipAction(5, someTime)

		expectedCosts := []ShipActionCost{
			{
				Resource: s.Costs[0].Resource,
				Amount:   180,
			},
			{
				Resource: s.Costs[1].Resource,
				Amount:   75,
			},
			{
				Resource: s.Costs[2].Resource,
				Amount:   500,
			},
			{
				Resource: s.Costs[3].Resource,
				Amount:   750,
			},
		}
		assert.Equal(t, expectedCosts, action.Costs)
	})

	t.Run("correctly calculates completion time when no resource is used", func(t *testing.T) {
		s := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs:     []ShipCost{},
		}

		action := s.CreateShipAction(5, someTime)

		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime, action.NextCompletionAt)
		assert.Equal(t, someTime, action.CompletedAt)
	})

	t.Run("correctly calculates completion time when single resource is used", func(t *testing.T) {
		s := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              uuid.New(),
					Cost:                  36,
					BuildTimeHoursPerUnit: 0.0004,
				},
			},
		}

		action := s.CreateShipAction(5, someTime)

		assert.Equal(t, someTime, action.CreatedAt)
		completionTime := 51840 * time.Millisecond
		assert.Equal(t, someTime.Add(completionTime), action.NextCompletionAt)
		totalTime := 5 * completionTime
		assert.Equal(t, someTime.Add(totalTime), action.CompletedAt)
	})

	t.Run("correctly calculates completion time when resource has no build time", func(t *testing.T) {
		s := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              crystalResourceId,
					Cost:                  79,
					BuildTimeHoursPerUnit: 0.0,
				},
			},
		}

		action := s.CreateShipAction(5, someTime)

		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime, action.NextCompletionAt)
		assert.Equal(t, someTime, action.CompletedAt)
	})

	t.Run("correctly calculates completion time when multiple resources are used", func(t *testing.T) {
		s := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              uuid.New(),
					Cost:                  12,
					BuildTimeHoursPerUnit: 1,
				},
				{
					Resource:              uuid.New(),
					Cost:                  87,
					BuildTimeHoursPerUnit: 36,
				},
				{
					Resource:              uuid.New(),
					Cost:                  106,
					BuildTimeHoursPerUnit: 0.04,
				},
				{
					Resource:              uuid.New(),
					Cost:                  201,
					BuildTimeHoursPerUnit: 0,
				},
			},
		}

		action := s.CreateShipAction(5, someTime)

		assert.Equal(t, someTime, action.CreatedAt)
		completionTime := 11333664 * time.Second
		assert.Equal(t, someTime.Add(completionTime), action.NextCompletionAt)
		totalTime := 5 * completionTime
		assert.Equal(t, someTime.Add(totalTime), action.CompletedAt)
	})
}

func generateTestShip(
	t *testing.T,
	modifiers ...func(*testing.T, *Ship),
) Ship {
	t.Helper()

	s := Ship{
		Id:        shipId,
		Name:      "test-ship",
		CreatedAt: someTime,
		Costs:     []ShipCost{},
	}

	for _, modifier := range modifiers {
		modifier(t, &s)
	}

	return s
}

func withShipCost(t *testing.T, s *Ship) {
	t.Helper()

	s.Costs = []ShipCost{
		{
			Resource:              metalResourceId,
			Cost:                  36,
			BuildTimeHoursPerUnit: 0.0004,
		},
		{
			Resource:              crystalResourceId,
			Cost:                  78,
			BuildTimeHoursPerUnit: 0.1,
		},
	}
}
