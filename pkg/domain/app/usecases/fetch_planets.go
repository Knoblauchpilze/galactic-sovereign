package usecases

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type FetchPlanetsUseCase struct {
	planetRepo    drivenports.ForListingPlanets
	planetMutator drivenports.ForMutatingPlanet
	clock         drivenports.ForFetchingTime
}

func NewFetchPlanetsUseCase(
	planetRepo drivenports.ForListingPlanets,
	planetMutator drivenports.ForMutatingPlanet,
	clock drivenports.ForFetchingTime,
) *FetchPlanetsUseCase {
	return &FetchPlanetsUseCase{
		planetRepo:    planetRepo,
		planetMutator: planetMutator,
		clock:         clock,
	}
}

func (p *FetchPlanetsUseCase) Get(ctx context.Context, id uuid.UUID) (models.Planet, error) {
	moment := p.clock.Now(ctx)
	result, err := p.planetMutator.Mutate(ctx, id, generateUpdateMutator(moment))
	if err != nil {
		return models.Planet{}, err
	}

	return result, nil
}

func (p *FetchPlanetsUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	moment := p.clock.Now(ctx)

	ids, err := p.planetRepo.ListForPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	out := make([]models.Planet, 0, len(ids))

	for _, id := range ids {
		result, err := p.planetMutator.Mutate(ctx, id, generateUpdateMutator(moment))
		if err != nil {
			return nil, err
		}

		out = append(out, result)
	}

	return out, nil
}

func generateUpdateMutator(moment time.Time) drivenports.PlanetMutator {
	return func(p *models.Planet) error {
		return domainservices.AdvancePlanetToTime(p, moment)
	}
}
