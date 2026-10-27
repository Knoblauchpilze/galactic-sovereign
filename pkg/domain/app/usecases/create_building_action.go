package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
)

type CreateBuildingActionUseCase struct {
	buildingRepo drivenports.ForFetchingBuilding
	simulator    drivenports.ForExecutingSimulationStep
	clock        drivenports.ForFetchingTime
}

func NewCreateBuildingActionUseCase(
	buildingRepo drivenports.ForFetchingBuilding,
	simulator drivenports.ForExecutingSimulationStep,
	clock drivenports.ForFetchingTime,
) *CreateBuildingActionUseCase {
	return &CreateBuildingActionUseCase{
		buildingRepo: buildingRepo,
		simulator:    simulator,
		clock:        clock,
	}
}

func (b *CreateBuildingActionUseCase) Create(
	ctx context.Context,
	req request.BuildingActionCreationRequest,
) (models.BuildingAction, error) {
	moment := b.clock.Now(ctx)

	building, err := b.buildingRepo.Get(ctx, req.Building)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return models.BuildingAction{}, domainerrors.ErrBuildingNotFound
		}

		return models.BuildingAction{}, err
	}

	mutator := generateBuildingActionMutator(building)
	result, err := domainservices.AdvanceSimulationUntil(ctx, b.simulator, req.Planet, moment, mutator)
	if err != nil {
		return models.BuildingAction{}, err
	}

	if result.Planet == nil {
		return models.BuildingAction{}, domainerrors.ErrNotFound
	}

	if result.Planet.BuildingAction == nil {
		return models.BuildingAction{}, domainerrors.ErrResourceCreationFailed
	}

	return *result.Planet.BuildingAction, nil
}

func generateBuildingActionMutator(building models.Building) models.SimulationStepResolver {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		err := step.Planet.AddBuildingAction(building)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		return models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     step.Planet,
		}, nil
	}
}
