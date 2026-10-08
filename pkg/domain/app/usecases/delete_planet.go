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
	deleter drivenports.ForDeletingPlanet
	clock   drivenports.ForFetchingTime
}

func NewDeletePlanetUseCase(
	deleter drivenports.ForDeletingPlanet,
	clock drivenports.ForFetchingTime,
) *DeletePlanetUseCase {
	return &DeletePlanetUseCase{
		deleter: deleter,
		clock:   clock,
	}
}

func (p *DeletePlanetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	moment := p.clock.Now(ctx)

	deleter := domainservices.AdvancePlanetToTimeThen(moment, planetDeletionGuard)
	err := p.deleter.Delete(ctx, id, deleter)
	if err != nil {
		return err
	}

	return nil
}

func planetDeletionGuard(p *models.Planet) error {
	if p.Homeworld {
		return domainerrors.ErrHomeworldCannotBeDeleted
	}

	if p.BuildingAction != nil {
		return domainerrors.ErrBuildingActionNotCompleted
	}

	if len(p.ShipActions) > 0 {
		return domainerrors.ErrShipActionNotCompleted
	}

	return nil
}
