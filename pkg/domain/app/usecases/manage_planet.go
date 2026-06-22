package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
)

type planetUseCase struct {
	playerRepo   drivenports.ForManagingPlayers
	universeRepo drivenports.ForManagingUniverses
	planetRepo   drivenports.ForManagingPlanets
}

func NewPlanetUseCase(
	playerRepo drivenports.ForManagingPlayers,
	universeRepo drivenports.ForManagingUniverses,
	planetRepo drivenports.ForManagingPlanets,
) drivingports.ForManagingPlanet {
	return &planetUseCase{
		playerRepo:   playerRepo,
		universeRepo: universeRepo,
		planetRepo:   planetRepo,
	}
}

func (p *planetUseCase) Create(ctx context.Context, req request.PlanetCreationRequest) (models.Planet, error) {
	player, err := p.playerRepo.Get(ctx, req.Player)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return models.Planet{}, domainerrors.ErrPlayerNotFound
		}
		return models.Planet{}, err
	}

	universe, err := p.universeRepo.Get(ctx, player.Universe)
	if err != nil {
		return models.Planet{}, err
	}

	planet := player.Colonize(universe)

	err = p.planetRepo.Create(ctx, planet)
	if err != nil {
		return models.Planet{}, err
	}

	return planet, nil
}

func (p *planetUseCase) Get(ctx context.Context, id uuid.UUID) (models.Planet, error) {
	return p.planetRepo.Get(ctx, id)
}

func (p *planetUseCase) List(ctx context.Context) ([]models.Planet, error) {
	return p.planetRepo.List(ctx)
}

func (p *planetUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	return p.planetRepo.ListForPlayer(ctx, player)
}

func (p *planetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return p.planetRepo.Delete(ctx, id)
}
