package usecases

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type DeleteBuildingActionUseCase struct {
	planetMutator drivenports.ForMutatingPlanet
	clock         drivenports.ForFetchingTime
}

func NewDeleteBuildingActionUseCase(
	planetMutator drivenports.ForMutatingPlanet,
	clock drivenports.ForFetchingTime,
) *DeleteBuildingActionUseCase {
	return &DeleteBuildingActionUseCase{
		planetMutator: planetMutator,
		clock:         clock,
	}
}

func (b *DeleteBuildingActionUseCase) DeleteForPlanet(
	ctx context.Context,
	planet uuid.UUID,
) error {
	moment := b.clock.Now(ctx)

	mutator := generateActionDeletionMutator(moment)
	result, err := b.planetMutator.Mutate(ctx, planet, mutator)
	if err != nil {
		return err
	}
	if result.Deleted {
		return domainerrors.ErrNotFound
	}

	return nil
}

func generateActionDeletionMutator(moment time.Time) drivenports.PlanetMutator {
	return func(p *models.Planet) (bool, error) {
		err := domainservices.AdvancePlanetToTime(p, moment)
		if err != nil {
			return false, err
		}

		p.BuildingAction = nil
		p.Version++

		return false, nil
	}
}
