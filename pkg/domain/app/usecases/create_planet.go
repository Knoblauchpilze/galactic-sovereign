package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
)

type CreatePlanetUseCase struct {
	playerRepo   drivenports.ForManagingPlayers
	universeRepo drivenports.ForManagingUniverses
	planetRepo   drivenports.ForCreatingPlanets
}

func NewCreatePlanetUseCase(
	playerRepo drivenports.ForManagingPlayers,
	universeRepo drivenports.ForManagingUniverses,
	planetRepo drivenports.ForCreatingPlanets,
) *CreatePlanetUseCase {
	return &CreatePlanetUseCase{
		playerRepo:   playerRepo,
		universeRepo: universeRepo,
		planetRepo:   planetRepo,
	}
}

func (p *CreatePlanetUseCase) Create(ctx context.Context, req request.PlanetCreationRequest) (models.Planet, error) {
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
