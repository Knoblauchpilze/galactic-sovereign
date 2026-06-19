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

type playerUseCase struct {
	playerRepo   drivenports.ForManagingPlayers
	resourceRepo drivenports.ForListingResources
	buildingRepo drivenports.ForListingBuildings
	planetRepo   drivenports.ForManagingPlanets
}

func NewPlayerUseCase(
	playerRepo drivenports.ForManagingPlayers,
	resourceRepo drivenports.ForListingResources,
	buildingRepo drivenports.ForListingBuildings,
	planetRepo drivenports.ForManagingPlanets,
) drivingports.ForManagingPlayer {
	return &playerUseCase{
		playerRepo:   playerRepo,
		resourceRepo: resourceRepo,
		buildingRepo: buildingRepo,
		planetRepo:   planetRepo,
	}
}

func (p *playerUseCase) Create(ctx context.Context, req request.PlayerCreationRequest) (models.Player, error) {
	player := request.FromPlayerCreationRequest(req)

	resources, err := p.resourceRepo.List(ctx)
	if err != nil {
		return models.Player{}, err
	}

	buildings, err := p.buildingRepo.List(ctx)
	if err != nil {
		return models.Player{}, err
	}

	homeworld := player.CreateHomeworld(resources, buildings)

	err = p.playerRepo.Create(ctx, player, homeworld)
	if err != nil {
		return models.Player{}, err
	}

	return player, nil
}

func (p *playerUseCase) Get(ctx context.Context, id uuid.UUID) (models.Player, error) {
	return p.playerRepo.Get(ctx, id)
}

func (p *playerUseCase) List(ctx context.Context) ([]models.Player, error) {
	return p.playerRepo.List(ctx)
}

func (p *playerUseCase) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]models.Player, error) {
	return p.playerRepo.ListForApiUser(ctx, apiUser)
}

func (p *playerUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	player, err := p.playerRepo.Get(ctx, id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return nil
		}

		return err
	}

	err = p.playerRepo.Delete(ctx, player)
	if err != nil {
		return err
	}

	return nil
}
