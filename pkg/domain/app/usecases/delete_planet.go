package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
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

	deleter := func(p *models.Planet) error {
		return domainservices.PlanetDeletionGuard(p, moment)
	}

	err := p.deleter.Delete(ctx, id, deleter)
	if err != nil {
		return err
	}

	return nil
}
