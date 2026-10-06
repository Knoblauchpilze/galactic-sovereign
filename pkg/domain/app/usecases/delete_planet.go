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

type DeletePlanetUseCase struct {
	planetMutator drivenports.ForMutatingPlanet
	clock         drivenports.ForFetchingTime
}

func NewDeletePlanetUseCase(
	planetMutator drivenports.ForMutatingPlanet,
	clock drivenports.ForFetchingTime,
) *DeletePlanetUseCase {
	return &DeletePlanetUseCase{
		planetMutator: planetMutator,
		clock:         clock,
	}
}

func (p *DeletePlanetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
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
			return false, domainerrors.ErrBuildingActionNotCompleted
		}

		if len(p.ShipActions) > 0 {
			return false, domainerrors.ErrShipActionNotCompleted
		}

		return true, nil
	}
}
