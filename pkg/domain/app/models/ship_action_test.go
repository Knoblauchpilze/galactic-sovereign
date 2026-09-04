package models

import (
	"testing"
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_ShipAction_CompleteOne(t *testing.T) {
	t.Run("returns error when all units are already completed", func(t *testing.T) {
		action := generateTestShipAction(t)
		action.Count = 0

		err := action.CompleteOne()

		assert.ErrorIs(t, err, domainerrors.ErrShipActionAlreadyCompleted, "Actual err: %v", err)
	})

	t.Run("decreases count by one", func(t *testing.T) {
		action := generateTestShipAction(t)
		action.Count = 2

		err := action.CompleteOne()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, 1, action.Count)
	})

	t.Run("updates next completion time based on unit completion time", func(t *testing.T) {
		action := generateTestShipAction(t)
		originalCompletionTime := action.NextCompletionAt
		action.Count = 2

		err := action.CompleteOne()
		require.NoError(t, err, "Actual err: %v", err)

		expected := originalCompletionTime.Add(1 * time.Hour)
		assert.Equal(t, expected, action.NextCompletionAt)
	})
}

func TestUnit_ShipAction_Completed(t *testing.T) {
	t.Run("returns true when count is 0", func(t *testing.T) {
		action := ShipAction{
			Id:    uuid.New(),
			Count: 0,
		}

		assert.True(t, action.Completed())
	})

	t.Run("returns true when count is negative", func(t *testing.T) {
		action := ShipAction{
			Id:    uuid.New(),
			Count: -1,
		}

		assert.True(t, action.Completed())
	})

	t.Run("returns false when count is positive", func(t *testing.T) {
		action := ShipAction{
			Id:    uuid.New(),
			Count: 12,
		}

		assert.False(t, action.Completed())
	})
}

func TestUnit_ShipAction_CompletionTime(t *testing.T) {
	t.Run("returns next completion time when only one unit is left", func(t *testing.T) {
		action := generateTestShipAction(t)
		action.Count = 1

		actual := action.CompletionTime()

		assert.Equal(t, action.NextCompletionAt, actual)
	})

	t.Run("returns next completion time when no unit is left", func(t *testing.T) {
		action := generateTestShipAction(t)
		action.Count = 0

		actual := action.CompletionTime()

		assert.Equal(t, action.NextCompletionAt, actual)
	})

	t.Run("returns correct completion time when multiple units are left", func(t *testing.T) {
		action := generateTestShipAction(t)
		require.NotEqual(t, 1, action.Count)

		actual := action.CompletionTime()

		completionTime := time.Duration(action.Count) * action.UnitCompletionTime
		expected := action.CreatedAt.Add(completionTime)
		assert.Equal(t, expected, actual)
	})

	t.Run("returns correct completion time when some units have already been produced", func(t *testing.T) {
		action := ShipAction{
			Id:        uuid.New(),
			Ship:      lightFighterId,
			Count:     10,
			CreatedAt: someTime,
			// This materializes that 2 units have already been produced
			NextCompletionAt:   someTime.Add(3 * time.Hour),
			UnitCompletionTime: 1 * time.Hour,
			Costs:              []ShipActionCost{},
		}

		actual := action.CompletionTime()

		completionTime := time.Duration(action.Count-1) * action.UnitCompletionTime
		expected := action.NextCompletionAt.Add(completionTime)
		assert.Equal(t, expected, actual)
	})
}

func generateTestShipAction(
	t *testing.T,
	modifiers ...func(*testing.T, *ShipAction),
) ShipAction {
	t.Helper()

	action := ShipAction{
		Id:                 uuid.New(),
		Ship:               lightFighterId,
		Count:              12,
		CreatedAt:          someTime,
		NextCompletionAt:   someTime.Add(1 * time.Hour),
		UnitCompletionTime: 1 * time.Hour,
		Costs:              []ShipActionCost{},
	}

	for _, modifier := range modifiers {
		modifier(t, &action)
	}

	return action
}
