package usecases

import (
	"context"

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

	mutator := func(p *models.Planet) error {
		return domainservices.AdvancePlanetToTime(ctx, p, moment)
	}

	return p.planetMutator.Mutate(ctx, id, mutator)
}

// TODO: Should make the planet up to date and save it
func (p *planetUseCase) List(ctx context.Context) ([]models.Planet, error) {
	return p.planetRepo.List(ctx)
}

// TODO: Should make the planet up to date and save it
func (p *planetUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	return p.planetRepo.ListForPlayer(ctx, player)
}

// TODO: Should make the planet up to date and save it
func (p *planetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return p.planetRepo.Delete(ctx, id)
}
