package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type FetchPlanetsUseCase struct {
	planetRepo drivenports.ForListingPlanets
	simulator  drivenports.ForExecutingSimulationStep
	clock      drivenports.ForFetchingTime
}

func NewFetchPlanetsUseCase(
	planetRepo drivenports.ForListingPlanets,
	simulator drivenports.ForExecutingSimulationStep,
	clock drivenports.ForFetchingTime,
) *FetchPlanetsUseCase {
	return &FetchPlanetsUseCase{
		planetRepo: planetRepo,
		simulator:  simulator,
		clock:      clock,
	}
}

func (p *FetchPlanetsUseCase) Get(ctx context.Context, id uuid.UUID) (models.Planet, error) {
	moment := p.clock.Now(ctx)

	result, err := domainservices.AdvanceSimulationUntil(ctx, p.simulator, id, moment, noOpMutator)
	if err != nil {
		return models.Planet{}, err
	}

	if result.Planet == nil {
		return models.Planet{}, domainerrors.ErrNotFound
	}

	return *result.Planet, nil
}

func (p *FetchPlanetsUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	moment := p.clock.Now(ctx)

	ids, err := p.planetRepo.ListForPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	out := make([]models.Planet, 0, len(ids))

	for _, id := range ids {

		result, err := domainservices.AdvanceSimulationUntil(ctx, p.simulator, id, moment, noOpMutator)
		if err != nil {
			return nil, err
		}

		if result.Planet != nil {
			out = append(out, *result.Planet)
		}
	}

	return out, nil
}

func noOpMutator(step models.SimulationStep) (models.SimulationStepResult, error) {
	return models.SimulationStepResult{
		Resolution: models.PlanetUpdated,
		Planet:     step.Planet,
	}, nil
}
