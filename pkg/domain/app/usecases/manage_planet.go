package usecases

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type planetUseCase struct {
	planetRepo    drivenports.ForManagingPlanets
	planetMutator drivenports.ForMutatingPlanet
	clock         drivenports.ForFetchingTime
}

func NewPlanetUseCase(
	planetRepo drivenports.ForManagingPlanets,
	planetMutator drivenports.ForMutatingPlanet,
	clock drivenports.ForFetchingTime,
) drivingports.ForManagingPlanet {
	return &planetUseCase{
		planetRepo:    planetRepo,
		planetMutator: planetMutator,
		clock:         clock,
	}
}

func (p *planetUseCase) Get(ctx context.Context, id uuid.UUID) (models.Planet, error) {
	moment := p.clock.Now(ctx)
	return p.planetMutator.Mutate(ctx, id, generateUpdateMutator(ctx, moment))
}

// TODO: Should make the planet up to date and save it
func (p *planetUseCase) List(ctx context.Context) ([]models.Planet, error) {
	return p.planetRepo.List(ctx)
}

func (p *planetUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	moment := p.clock.Now(ctx)

	ids, err := p.planetRepo.ListForPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	out := make([]models.Planet, 0, len(ids))

	for _, id := range ids {
		planet, err := p.planetMutator.Mutate(ctx, id, generateUpdateMutator(ctx, moment))
		if err != nil {
			return nil, err
		}

		out = append(out, planet)
	}

	return out, nil
}

// TODO: Should make the planet up to date and save it
// TODO: It is not needed to update: even when points are introduced, the
// points will be added and immediately deleted. It makes sense to strenghten
// the delete though so that a building action cannot be running when a planet
// is deleted.
func (p *planetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return p.planetRepo.Delete(ctx, id)
}

func generateUpdateMutator(
	ctx context.Context,
	moment time.Time,
) drivenports.PlanetMutator {
	return func(p *models.Planet) error {
		return domainservices.AdvancePlanetToTime(ctx, p, moment)
	}
}
