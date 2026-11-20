package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type FetchFleetsUseCase struct {
	playerRepo drivenports.ForManagingPlayers
	simulator  drivenports.ForExecutingSimulationStep
	clock      drivenports.ForFetchingTime
}

func NewFetchFleetsUseCase(
	playerRepo drivenports.ForManagingPlayers,
	simulator drivenports.ForExecutingSimulationStep,
	clock drivenports.ForFetchingTime,
) *FetchFleetsUseCase {
	return &FetchFleetsUseCase{
		playerRepo: playerRepo,
		simulator:  simulator,
		clock:      clock,
	}
}

func (u *FetchFleetsUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Fleet, error) {
	playerData, err := u.playerRepo.Get(ctx, player)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return []models.Fleet{}, domainerrors.ErrPlayerNotFound
		}

		return []models.Fleet{}, err
	}

	moment := u.clock.Now(ctx)

	err = domainservices.AdvanceSimulationForUniverseUntil(ctx, u.simulator, playerData.Universe, moment)
	if err != nil {
		return []models.Fleet{}, err
	}

	// TODO: Explain that it's acceptable to update the simulation and then read the fleets.
	// The delay between both operation is small and there's no real way to make it work in
	// an atomic way without acquiring some sort of 'lock' on the whole universe.

	return []models.Fleet{}, nil
}
