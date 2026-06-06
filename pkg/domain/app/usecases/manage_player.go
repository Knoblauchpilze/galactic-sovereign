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
	repo drivenports.ForManagingPlayers
}

func NewPlayerUseCase(repo drivenports.ForManagingPlayers) drivingports.ForManagingPlayer {
	return &playerUseCase{
		repo: repo,
	}
}

func (p *playerUseCase) Create(ctx context.Context, req request.PlayerCreationRequest) (models.Player, error) {
	player := request.FromPlayerCreationRequest(req)

	err := p.repo.Create(ctx, player)
	if err != nil {
		return models.Player{}, err
	}

	return player, nil
}

func (p *playerUseCase) Get(ctx context.Context, id uuid.UUID) (models.Player, error) {
	return p.repo.Get(ctx, id)
}

func (p *playerUseCase) List(ctx context.Context) ([]models.Player, error) {
	return p.repo.List(ctx)
}

func (p *playerUseCase) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]models.Player, error) {
	return p.repo.ListForApiUser(ctx, apiUser)
}

func (p *playerUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	err := p.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
