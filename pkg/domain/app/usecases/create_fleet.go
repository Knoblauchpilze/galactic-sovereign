package usecases

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
)

type CreateFleetUseCase struct {
	universeRepo drivenports.ForFetchingUniverses
	dispatcher   drivenports.ForDispatchingFleet
	clock        drivenports.ForFetchingTime
}

func NewCreateFleetUseCase(
	universeRepo drivenports.ForFetchingUniverses,
	dispatcher drivenports.ForDispatchingFleet,
	clock drivenports.ForFetchingTime,
) *CreateFleetUseCase {
	return &CreateFleetUseCase{
		universeRepo: universeRepo,
		dispatcher:   dispatcher,
		clock:        clock,
	}
}

func (f *CreateFleetUseCase) Create(
	ctx context.Context,
	req request.FleetCreationRequest,
) (models.Fleet, error) {
	u, err := f.universeRepo.GetByPlanetId(ctx, req.Planet)
	if err != nil {
		return models.Fleet{}, err
	}

	coord := req.Destination.ToCoordinates()
	if !u.ValidCoordinates(coord) {
		return models.Fleet{}, domainerrors.ErrCoordinatesOutOfBound
	}

	moment := f.clock.Now(ctx)

	mutation := generateFleetMutation(moment)
	fleet, err := f.dispatcher.Dispatch(ctx, req, mutation)
	if err != nil {
		return models.Fleet{}, err
	}

	return fleet, nil

}

func generateFleetMutation(
	moment time.Time,
) drivenports.FleetCreator {
	return func(p *models.Planet, _ request.FleetCreationRequest) (models.Fleet, error) {
		err := domainservices.AdvancePlanetToTime(p, moment)
		if err != nil {
			return models.Fleet{}, err
		}

		return p.CreateFleet()
	}
}
