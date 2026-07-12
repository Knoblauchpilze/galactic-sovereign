package usecases

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type PlanetUseCase struct {
	planetRepo    drivenports.ForManagingPlanets
	planetMutator drivenports.ForMutatingPlanet
	clock         drivenports.ForFetchingTime
}

func NewPlanetUseCase(
	planetRepo drivenports.ForManagingPlanets,
	planetMutator drivenports.ForMutatingPlanet,
	clock drivenports.ForFetchingTime,
) *PlanetUseCase {
	return &PlanetUseCase{
		planetRepo:    planetRepo,
		planetMutator: planetMutator,
		clock:         clock,
	}
}

func (p *PlanetUseCase) Get(ctx context.Context, id uuid.UUID) (models.Planet, error) {
	moment := p.clock.Now(ctx)
	result, err := p.planetMutator.Mutate(ctx, id, generateUpdateMutator(moment))
	if err != nil {
		return models.Planet{}, err
	}

	if result.Deleted {
		return models.Planet{}, domainerrors.ErrNotFound
	}

	return result.Planet, nil
}

func (p *PlanetUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	moment := p.clock.Now(ctx)

	ids, err := p.planetRepo.ListForPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	out := make([]models.Planet, 0, len(ids))

	for _, id := range ids {
		result, err := p.planetMutator.Mutate(ctx, id, generateUpdateMutator(moment))
		if err != nil {
			return nil, err
		}

		if !result.Deleted {
			out = append(out, result.Planet)
		}
	}

	return out, nil
}

func (p *PlanetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	moment := p.clock.Now(ctx)

	result, err := p.planetMutator.Mutate(ctx, id, generateDeleteMutator(moment))
	if err != nil {
		return err
	}

	if !result.Deleted {
		return domainerrors.ErrPlanetDeletionFailed
	}

	return nil
}

func generateUpdateMutator(moment time.Time) drivenports.PlanetMutator {
	return func(p *models.Planet) (bool, error) {
		return false, domainservices.AdvancePlanetToTime(p, moment)
	}
}

func generateDeleteMutator(moment time.Time) drivenports.PlanetMutator {
	return func(p *models.Planet) (bool, error) {
		err := domainservices.AdvancePlanetToTime(p, moment)
		if err != nil {
			return false, err
		}

		if p.Homeworld {
			return false, domainerrors.ErrHomeworldCannotBeDeleted
		}

		if p.BuildingAction != nil {
			return false, domainerrors.ErrActionNotCompleted
		}

		return true, nil
	}
}
