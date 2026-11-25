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
	fleetRepo  drivenports.ForListingFleets
	clock      drivenports.ForFetchingTime
}

func NewFetchFleetsUseCase(
	playerRepo drivenports.ForManagingPlayers,
	simulator drivenports.ForExecutingSimulationStep,
	fleetRepo drivenports.ForListingFleets,
	clock drivenports.ForFetchingTime,
) *FetchFleetsUseCase {
	return &FetchFleetsUseCase{
		playerRepo: playerRepo,
		simulator:  simulator,
		fleetRepo:  fleetRepo,
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

	// It's acceptable to update the simulation first and then read the fleets in another
	// transaction. The delay between both operations is small and this flow only reads so
	// it does not modify the state of fleets.
	// A way to make it atomic would be to extend the simulator to allow performing not only
	// an update of a fleet or planet but also to fetch additional data. This seems out of
	// proportion for the use cases available at the moment.
	fleets, err := u.fleetRepo.ListForPlayer(ctx, player)
	if err != nil {
		return []models.Fleet{}, err
	}

	return fleets, nil
}
