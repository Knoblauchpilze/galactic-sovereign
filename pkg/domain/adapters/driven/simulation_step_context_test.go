package drivenadapters

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_SimulationStepContext_Resolve(t *testing.T) {
	t.Run("returns error when resolver requests planet deletion but is not allowed", func(t *testing.T) {
		ctx := simulationStepContext{
			step: models.SimulationStep{
				Timestamp: someTime,
				Planet: &models.Planet{
					Id: uuid.New(),
				},
				NextFleetEvent: nil,
			},
			allowPlanetDeletion: false,
			resolver: func(_ models.SimulationStep) (models.SimulationStepResult, error) {
				return models.SimulationStepResult{
					Resolution: models.PlanetDeleted,
					Planet:     nil,
					Fleet:      nil,
				}, nil
			},
		}

		err := ctx.resolve()

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})

	t.Run("accepts deletion when resolver requests planet deletion", func(t *testing.T) {
		ctx := simulationStepContext{
			step: models.SimulationStep{
				Timestamp: someTime,
				Planet: &models.Planet{
					Id: uuid.New(),
				},
				NextFleetEvent: nil,
			},
			allowPlanetDeletion: true,
			resolver: func(_ models.SimulationStep) (models.SimulationStepResult, error) {
				return models.SimulationStepResult{
					Resolution: models.PlanetDeleted,
					Planet:     nil,
					Fleet:      nil,
				}, nil
			},
		}

		err := ctx.resolve()
		require.NoError(t, err, "Actual err: %v", err)
	})
}
