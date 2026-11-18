package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

type PlayerUseCase struct {
	playerRepo   drivenports.ForManagingPlayers
	universeRepo drivenports.ForFetchingUniverses
}

func NewPlayerUseCase(
	playerRepo drivenports.ForManagingPlayers,
	universeRepo drivenports.ForFetchingUniverses,
) *PlayerUseCase {
	return &PlayerUseCase{
		playerRepo:   playerRepo,
		universeRepo: universeRepo,
	}
}

func (u *PlayerUseCase) Create(
	ctx context.Context,
	req request.PlayerCreationRequest,
) (models.Player, error) {
	player := request.FromPlayerCreationRequest(req)

	universe, err := u.universeRepo.Get(ctx, player.Universe)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return models.Player{}, domainerrors.ErrUniverseNotFound
		}
		return models.Player{}, err
	}

	homeworld := player.CreateHomeworld(universe)

	// TODO: Could be that the planet was used, in which case an optimistic lock error will be returned
	err = u.playerRepo.Create(ctx, player, homeworld)
	if err != nil {
		return models.Player{}, err
	}

	return player, nil
}

func (u *PlayerUseCase) Get(ctx context.Context, id uuid.UUID) (models.Player, error) {
	return u.playerRepo.Get(ctx, id)
}

func (u *PlayerUseCase) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]models.Player, error) {
	return u.playerRepo.ListForApiUser(ctx, apiUser)
}

func (u *PlayerUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: The simulation step should probably be ran to get the updated planet count
	player, err := u.playerRepo.Get(ctx, id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return nil
		}

		return err
	}

	err = u.playerRepo.Delete(ctx, player)
	if err != nil {
		return err
	}

	return nil
}
