package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
)

type planetUseCase struct {
	repo drivenports.ForManagingPlanets
}

func NewPlanetUseCase(repo drivenports.ForManagingPlanets) drivingports.ForManagingPlanet {
	return &planetUseCase{
		repo: repo,
	}
}

func (p *planetUseCase) Create(ctx context.Context, req request.PlanetCreationRequest) (models.Planet, error) {
	planet := request.FromPlanetCreationRequest(req)

	err := p.repo.Create(ctx, planet)
	if err != nil {
		return models.Planet{}, err
	}

	return planet, nil
}

func (p *planetUseCase) Get(ctx context.Context, id uuid.UUID) (models.Planet, error) {
	return p.repo.Get(ctx, id)
}

func (p *planetUseCase) List(ctx context.Context) ([]models.Planet, error) {
	return p.repo.List(ctx)
}

func (p *planetUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	return p.repo.ListForPlayer(ctx, player)
}

func (p *planetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return p.repo.Delete(ctx, id)
}
