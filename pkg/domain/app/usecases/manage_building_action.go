package usecases

import (
	"context"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
)

type buildingActionUseCase struct {
	actionRepo drivenports.ForManagingBuildingActions
	planetRepo drivenports.ForManagingPlanets
}

func NewBuildingActionUseCase(
	actionRepo drivenports.ForManagingBuildingActions,
	planetRepo drivenports.ForManagingPlanets,
) drivingports.ForManagingBuildingAction {
	return &buildingActionUseCase{
		actionRepo: actionRepo,
		planetRepo: planetRepo,
	}
}

// TODO: Should make the planet up to date and save it
func (b *buildingActionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	planet, err := b.planetRepo.GetByAction(ctx, id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return nil
		}
		return err
	}

	err = planet.CancelBuildingAction()
	if err != nil {
		if err == domainerrors.ErrNoActionInProgress {
			return nil
		}
		return err
	}

	err = b.actionRepo.Delete(ctx, planet)
	if err != nil {
		return err
	}

	return nil
}
