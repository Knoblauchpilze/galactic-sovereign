package usecases

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
)

type CreateBuildingActionUseCase struct {
	buildingRepo  drivenports.ForFetchingBuilding
	planetMutator drivenports.ForMutatingPlanet
	clock         drivenports.ForFetchingTime
}

func NewCreateBuildingActionUseCase(
	buildingRepo drivenports.ForFetchingBuilding,
	planetMutator drivenports.ForMutatingPlanet,
	clock drivenports.ForFetchingTime,
) *CreateBuildingActionUseCase {
	return &CreateBuildingActionUseCase{
		buildingRepo:  buildingRepo,
		planetMutator: planetMutator,
		clock:         clock,
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

	mutator := generateActionMutator(moment, building)
	result, err := b.planetMutator.Mutate(ctx, req.Planet, mutator)
	if err != nil {
		return models.BuildingAction{}, err
	}
	if result.Deleted {
		return models.BuildingAction{}, domainerrors.ErrNotFound
	}

	if result.Planet.BuildingAction == nil {
		return models.BuildingAction{}, domainerrors.ErrResourceCreationFailed
	}

	return *result.Planet.BuildingAction, nil
}

func generateActionMutator(moment time.Time, building models.Building) drivenports.PlanetMutator {
	return func(p *models.Planet) (bool, error) {
		err := domainservices.AdvancePlanetToTime(p, moment)
		if err != nil {
			return false, err
		}

		return false, p.AddBuildingAction(building)
	}
}
