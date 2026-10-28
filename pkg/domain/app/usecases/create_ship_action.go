package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
)

type CreateShipActionUseCase struct {
	shipRepo  drivenports.ForFetchingShip
	simulator drivenports.ForExecutingSimulationStep
	clock     drivenports.ForFetchingTime
}

func NewCreateShipActionUseCase(
	shipRepo drivenports.ForFetchingShip,
	simulator drivenports.ForExecutingSimulationStep,
	clock drivenports.ForFetchingTime,
) *CreateShipActionUseCase {
	return &CreateShipActionUseCase{
		shipRepo:  shipRepo,
		simulator: simulator,
		clock:     clock,
	}
}

func (u *CreateShipActionUseCase) Create(
	ctx context.Context,
	req request.ShipActionCreationRequest,
) (models.ShipAction, error) {
	moment := u.clock.Now(ctx)

	ship, err := u.shipRepo.Get(ctx, req.Ship)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return models.ShipAction{}, domainerrors.ErrShipNotFound
		}

		return models.ShipAction{}, err
	}

	mutator := generateShipActionMutator(ship, req.Count)
	result, err := domainservices.AdvanceSimulationUntil(ctx, u.simulator, req.Planet, moment, mutator)
	if err != nil {
		return models.ShipAction{}, err
	}

	if result.Planet == nil {
		return models.ShipAction{}, domainerrors.ErrNotFound
	}

	if len(result.Planet.ShipActions) == 0 {
		return models.ShipAction{}, domainerrors.ErrResourceCreationFailed
	}

	last := result.Planet.ShipActions[len(result.Planet.ShipActions)-1]
	return last, nil
}

func generateShipActionMutator(
	ship models.Ship,
	count int,
) models.SimulationStepResolver {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		err := step.Planet.AddShipAction(ship, count)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		return models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     step.Planet,
		}, nil
	}
}
