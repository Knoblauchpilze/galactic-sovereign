package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

type FetchSolarSystemUseCase struct {
	solarSystemRepo drivenports.ForFetchingSolarSystem
}

func NewFetchSolarSystemUseCase(
	solarSystemRepo drivenports.ForFetchingSolarSystem,
) *FetchSolarSystemUseCase {
	return &FetchSolarSystemUseCase{
		solarSystemRepo: solarSystemRepo,
	}
}

func (u *FetchSolarSystemUseCase) GetSolarSystem(
	ctx context.Context,
	universe uuid.UUID,
	galaxy int,
	solarSystem int,
) (models.SolarSystem, error) {
	// TODO: This might also need to be ran through the simulation step to process
	// colonizing fleets.
	return u.solarSystemRepo.GetSolarSystem(ctx, universe, galaxy, solarSystem)
}
