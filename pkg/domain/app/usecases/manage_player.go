package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
)

type playerUseCase struct {
	playerRepo   drivenports.ForManagingPlayers
	resourceRepo drivenports.ForListingResources
}

func NewPlayerUseCase(
	playerRepo drivenports.ForManagingPlayers,
	resourceRepo drivenports.ForListingResources,
) drivingports.ForManagingPlayer {
	return &playerUseCase{
		playerRepo:   playerRepo,
		resourceRepo: resourceRepo,
	}
}

func (p *playerUseCase) Create(ctx context.Context, req request.PlayerCreationRequest) (models.Player, error) {
	player := request.FromPlayerCreationRequest(req)
	homeworld := player.CreateHomeworld([]models.Resource{})

	err := p.playerRepo.Create(ctx, player, homeworld)
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
	err := p.playerRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
