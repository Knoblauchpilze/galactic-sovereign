package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type FetchSolarSystemUseCase struct {
	simulator       drivenports.ForExecutingSimulationStep
	solarSystemRepo drivenports.ForFetchingSolarSystem
	clock           drivenports.ForFetchingTime
}

func NewFetchSolarSystemUseCase(
	simulator drivenports.ForExecutingSimulationStep,
	solarSystemRepo drivenports.ForFetchingSolarSystem,
	clock drivenports.ForFetchingTime,
) *FetchSolarSystemUseCase {
	return &FetchSolarSystemUseCase{
		simulator:       simulator,
		solarSystemRepo: solarSystemRepo,
		clock:           clock,
	}
}

func (u *FetchSolarSystemUseCase) GetSolarSystem(
	ctx context.Context,
	universe uuid.UUID,
	galaxy int,
	solarSystem int,
) (models.SolarSystem, error) {
	moment := u.clock.Now(ctx)

	err := domainservices.AdvanceSimulationForUniverseUntil(ctx, u.simulator, universe, moment)
	if err != nil {
		return models.SolarSystem{}, err
	}

	// It's acceptable to update the simulation first and then read the solar system in another
	// transaction. The delay between both operations is small and this flow only reads so it
	// does not modify the state of fleets.
	// A way to make it atomic would be to extend the simulator to allow performing not only
	// an update of a fleet or planet but also to fetch additional data. This seems out of
	// proportion for the use cases available at the moment.
	// It's the second such use case besides fetching fleet.
	system, err := u.solarSystemRepo.GetSolarSystem(ctx, universe, galaxy, solarSystem)
	if err != nil {
		return models.SolarSystem{}, err
	}

	return system, nil
}
