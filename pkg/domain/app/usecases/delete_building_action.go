package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type DeleteBuildingActionUseCase struct {
	simulator drivenports.ForExecutingSimulationStep
	clock     drivenports.ForFetchingTime
}

func NewDeleteBuildingActionUseCase(
	simulator drivenports.ForExecutingSimulationStep,
	clock drivenports.ForFetchingTime,
) *DeleteBuildingActionUseCase {
	return &DeleteBuildingActionUseCase{
		simulator: simulator,
		clock:     clock,
	}
}

func (b *DeleteBuildingActionUseCase) DeleteForPlanet(
	ctx context.Context,
	planet uuid.UUID,
) error {
	moment := b.clock.Now(ctx)

	mutator := generateBuildingActionDeletionMutator()
	_, err := domainservices.AdvanceSimulationUntil(ctx, b.simulator, planet, moment, mutator)
	if err != nil {
		return err
	}

	return nil
}

func generateBuildingActionDeletionMutator() models.SimulationStepResolver {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		err := step.Planet.CancelBuildingAction()
		if err != nil {
			if err == domainerrors.ErrNoActionInProgress {
				return models.SimulationStepResult{Resolution: models.PlanetUpdated,
					Planet: step.Planet,
				}, nil
			}

			return models.SimulationStepResult{}, err
		}

		return models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     step.Planet,
		}, nil
	}
}
