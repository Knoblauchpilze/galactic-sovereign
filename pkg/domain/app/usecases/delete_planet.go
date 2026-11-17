package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type DeletePlanetUseCase struct {
	simulator drivenports.ForExecutingSimulationStep
	clock     drivenports.ForFetchingTime
}

func NewDeletePlanetUseCase(
	simulator drivenports.ForExecutingSimulationStep,
	clock drivenports.ForFetchingTime,
) *DeletePlanetUseCase {
	return &DeletePlanetUseCase{
		simulator: simulator,
		clock:     clock,
	}
}

func (p *DeletePlanetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	moment := p.clock.Now(ctx)

	result, err := domainservices.AdvanceSimulationForDeleteUntil(ctx, p.simulator, id, moment, planetDeletionGuard)
	if err != nil {
		return err
	}

	if result.Resolution != models.PlanetDeleted {
		return domainerrors.ErrInvalidSimulationStepResult
	}

	return nil
}

func planetDeletionGuard(step models.SimulationStep) (models.SimulationStepResult, error) {
	// func planetDeletionGuard(p *models.Planet) error {
	if step.Planet.Homeworld {
		return models.SimulationStepResult{}, domainerrors.ErrHomeworldCannotBeDeleted
	}

	if step.Planet.BuildingAction != nil {
		return models.SimulationStepResult{}, domainerrors.ErrBuildingActionNotCompleted
	}

	if len(step.Planet.ShipActions) > 0 {
		return models.SimulationStepResult{}, domainerrors.ErrShipActionNotCompleted
	}

	return models.SimulationStepResult{Resolution: models.PlanetDeleted}, nil
}
