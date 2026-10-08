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

	mutator := domainservices.AdvancePlanetToTimeThen(moment, actionDeletionMutator)
	_, err := b.planetMutator.Mutate(ctx, planet, mutator)
	if err != nil {
		return err
	}

	return nil
}

func actionDeletionMutator(p *models.Planet) error {
	err := p.CancelBuildingAction()
	if err != nil {
		if err == domainerrors.ErrNoActionInProgress {
			return nil
		}

		return err
	}

	return nil
}
