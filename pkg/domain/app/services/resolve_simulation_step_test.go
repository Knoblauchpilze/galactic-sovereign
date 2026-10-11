package domainservices

import (
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_ResolveSimulationStep(t *testing.T) {
	t.Run("advances planet to time when no fleet event is available", func(t *testing.T) {
		step := models.SimulationStep{
			Timestamp: t3,
			Planet: models.Planet{
				Id:        uuid.New(),
				CreatedAt: t1,
				UpdatedAt: t2,
				Version:   1,
			},
			NextFleetEvent: nil,
		}

		actual, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SimulationStepResult{
			Planet: &models.Planet{
				Id:        step.Planet.Id,
				CreatedAt: t1,
				UpdatedAt: t3,
				Version:   2,
			},
			Fleet: nil,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns unsupported when fleet event is available", func(t *testing.T) {
		step := models.SimulationStep{
			Timestamp: t3,
			Planet: models.Planet{
				Id:        uuid.New(),
				CreatedAt: t1,
				UpdatedAt: t2,
				Version:   1,
			},
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t2,
				Fleet:     models.Fleet{},
				Planet:    nil,
			},
		}

		_, err := ResolveSimulationStep(step)

		assert.ErrorIs(t, err, errors.ErrNotImplemented, "Actual err: %v", err)
	})
}
