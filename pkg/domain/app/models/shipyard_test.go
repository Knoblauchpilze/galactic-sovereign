package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Shipyard_CreateShipAction(t *testing.T) {
	t.Run("correctly calculates action costs", func(t *testing.T) {
		ship := generateTestShip(t, withShipCost)

		shipyard := Shipyard{throughput: 1.0}
		action := shipyard.CreateShipAction(ship, 5, someTime)

		expected := ShipAction{
			// The identifier is generated
			Id:        action.Id,
			Ship:      ship.Id,
			Count:     5,
			CreatedAt: someTime,
			// Ignore the completion here, there are dedicated tests
			NextCompletionAt:   action.NextCompletionAt,
			UnitCompletionTime: action.UnitCompletionTime,
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
		ship := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              metalResourceId,
					Cost:                  36,
					BuildTimeHoursPerUnit: 1,
				},
				{
					Resource:              crystalResourceId,
					Cost:                  15,
					BuildTimeHoursPerUnit: 36,
				},
				{
					Resource:              lightFighterId,
					Cost:                  100,
					BuildTimeHoursPerUnit: 0.04,
				},
				{
					Resource:              shipId,
					Cost:                  150,
					BuildTimeHoursPerUnit: 0,
				},
			},
		}

		shipyard := Shipyard{throughput: 1.0}
		action := shipyard.CreateShipAction(ship, 5, someTime)

		expectedCosts := []ShipActionCost{
			{
				Resource: metalResourceId,
				Amount:   180,
			},
			{
				Resource: crystalResourceId,
				Amount:   75,
			},
			{
				Resource: lightFighterId,
				Amount:   500,
			},
			{
				Resource: shipId,
				Amount:   750,
			},
		}
		assert.Equal(t, expectedCosts, action.Costs)
	})

	t.Run("correctly calculates completion time when no resource is used", func(t *testing.T) {
		shipyard := Shipyard{throughput: 1.0}
		ship := generateTestShip(t)

		action := shipyard.CreateShipAction(ship, 5, someTime)

		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime, action.NextCompletionAt)
		assert.Equal(t, time.Duration(0), action.UnitCompletionTime)
	})

	t.Run("correctly calculates completion time when single resource is used", func(t *testing.T) {
		shipyard := Shipyard{throughput: 1.0}
		ship := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              metalResourceId,
					Cost:                  36,
					BuildTimeHoursPerUnit: 0.0004,
				},
			},
		}

		action := shipyard.CreateShipAction(ship, 5, someTime)

		completionTime := 51840 * time.Millisecond
		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime.Add(completionTime), action.NextCompletionAt)
		assert.Equal(t, completionTime, action.UnitCompletionTime)
	})

	t.Run("correctly calculates completion time when resource has no build time", func(t *testing.T) {
		shipyard := Shipyard{throughput: 1.0}
		ship := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              crystalResourceId,
					Cost:                  79,
					BuildTimeHoursPerUnit: 0,
				},
			},
		}

		action := shipyard.CreateShipAction(ship, 5, someTime)

		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime, action.NextCompletionAt)
		assert.Equal(t, time.Duration(0), action.UnitCompletionTime)
	})

	t.Run("correctly calculates completion time when multiple resources are used", func(t *testing.T) {
		shipyard := Shipyard{throughput: 1.0}
		ship := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              metalResourceId,
					Cost:                  12,
					BuildTimeHoursPerUnit: 1,
				},
				{
					Resource:              crystalResourceId,
					Cost:                  87,
					BuildTimeHoursPerUnit: 36,
				},
				{
					Resource:              lightFighterId,
					Cost:                  106,
					BuildTimeHoursPerUnit: 0.04,
				},
				{
					Resource:              shipId,
					Cost:                  201,
					BuildTimeHoursPerUnit: 0,
				},
			},
		}

		action := shipyard.CreateShipAction(ship, 5, someTime)

		completionTime := 11333664 * time.Second
		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime.Add(completionTime), action.NextCompletionAt)
		assert.Equal(t, completionTime, action.UnitCompletionTime)
	})

	t.Run("correctly calculates completion time when throughput is greater than one", func(t *testing.T) {
		shipyard := Shipyard{throughput: 2.0}
		ship := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              metalResourceId,
					Cost:                  36,
					BuildTimeHoursPerUnit: 0.0004,
				},
			},
		}

		action := shipyard.CreateShipAction(ship, 5, someTime)

		completionTime := 25920 * time.Millisecond
		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime.Add(completionTime), action.NextCompletionAt)
		assert.Equal(t, completionTime, action.UnitCompletionTime)
	})

	t.Run("correctly applies throughput when multiple resources are used", func(t *testing.T) {
		shipyard := Shipyard{throughput: 1.05}
		ship := Ship{
			Id:        shipId,
			Name:      "test-ship",
			CreatedAt: someTime,
			Costs: []ShipCost{
				{
					Resource:              metalResourceId,
					Cost:                  12,
					BuildTimeHoursPerUnit: 1,
				},
				{
					Resource:              crystalResourceId,
					Cost:                  87,
					BuildTimeHoursPerUnit: 36,
				},
				{
					Resource:              lightFighterId,
					Cost:                  106,
					BuildTimeHoursPerUnit: 0.04,
				},
				{
					Resource:              shipId,
					Cost:                  201,
					BuildTimeHoursPerUnit: 0,
				},
			},
		}

		action := shipyard.CreateShipAction(ship, 5, someTime)

		completionTime := 10793965714285714 * time.Nanosecond
		assert.Equal(t, someTime, action.CreatedAt)
		assert.Equal(t, someTime.Add(completionTime), action.NextCompletionAt)
		assert.Equal(t, completionTime, action.UnitCompletionTime)
	})
}
