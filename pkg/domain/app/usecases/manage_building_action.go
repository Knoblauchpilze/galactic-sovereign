package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
)

type buildingActionUseCase struct {
	actionRepo   drivenports.ForManagingBuildingActions
	planetRepo   drivenports.ForManagingPlanets
	buildingRepo drivenports.ForListingBuildings
}

func NewBuildingActionUseCase(
	actionRepo drivenports.ForManagingBuildingActions,
	planetRepo drivenports.ForManagingPlanets,
	buildingRepo drivenports.ForListingBuildings,
) drivingports.ForManagingBuildingAction {
	return &buildingActionUseCase{
		actionRepo:   actionRepo,
		planetRepo:   planetRepo,
		buildingRepo: buildingRepo,
	}
}

func (b *buildingActionUseCase) Create(ctx context.Context, req request.BuildingActionCreationRequest) (models.BuildingAction, error) {
	planet, err := b.planetRepo.Get(ctx, req.Planet)
	if err != nil {
		return models.BuildingAction{}, err
	}

	building, err := b.buildingRepo.Get(ctx, req.Building)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return models.BuildingAction{}, domainerrors.ErrBuildingNotFound
		}

		return models.BuildingAction{}, err
	}

	action, err := planet.AddBuildingAction(building)
	if err != nil {
		return models.BuildingAction{}, err
	}

	err = b.actionRepo.Create(ctx, planet, action)
	if err != nil {
		return models.BuildingAction{}, err
	}

	return action, nil
}

func (b *buildingActionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	action, err := b.actionRepo.Get(ctx, id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return nil
		}

		return err
	}

	planet, err := b.planetRepo.Get(ctx, action.Planet)
	if err != nil {
		return err
	}

	err = b.actionRepo.Delete(ctx, planet, action)
	if err != nil {
		return err
	}

	return nil
}
