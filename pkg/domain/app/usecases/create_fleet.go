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
	planet uuid.UUID,
	order models.FleetOrder,
) (models.Fleet, error) {
	u, err := f.universeRepo.GetByPlanetId(ctx, planet)
	if err != nil {
		return models.Fleet{}, err
	}

	if !u.ValidCoordinates(order.Destination) {
		return models.Fleet{}, domainerrors.ErrCoordinatesOutOfBound
	}

	speed, err := domainservices.DetermineFleetSpeed(order, u.Ships)
	if err != nil {
		return models.Fleet{}, err
	}

	moment := f.clock.Now(ctx)

	mutation := generateFleetMutation(moment, order, speed, u.OccupancyMap)
	fleet, err := f.dispatcher.Dispatch(ctx, planet, mutation)
	if err != nil {
		return models.Fleet{}, err
	}

	return fleet, nil

}

func generateFleetMutation(
	moment time.Time,
	order models.FleetOrder,
	speed int,
	occupancy models.OccupancyMap,
) drivenports.FleetCreator {
	return func(p *models.Planet) (models.Fleet, error) {
		err := domainservices.AdvancePlanetToTime(p, moment)
		if err != nil {
			return models.Fleet{}, err
		}

		flight := models.FleetFlight{
			StartTime:   moment,
			Source:      p.Coordinate,
			Destination: order.Destination,
			Target:      nil,
			Speed:       speed,
		}

		if target, ok := occupancy.UsedSlots[order.Destination]; ok {
			flight.Target = &target
		}

		return p.CreateFleet(order, flight)
	}
}
