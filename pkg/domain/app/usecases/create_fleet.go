package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
	"github.com/google/uuid"
)

type CreateFleetUseCase struct {
	universeRepo drivenports.ForFetchingUniverses
	simulator    drivenports.ForExecutingSimulationStep
	clock        drivenports.ForFetchingTime
}

func NewCreateFleetUseCase(
	universeRepo drivenports.ForFetchingUniverses,
	simulator drivenports.ForExecutingSimulationStep,
	clock drivenports.ForFetchingTime,
) *CreateFleetUseCase {
	return &CreateFleetUseCase{
		universeRepo: universeRepo,
		simulator:    simulator,
		clock:        clock,
	}
}

func (u *CreateFleetUseCase) Create(
	ctx context.Context,
	planet uuid.UUID,
	order models.FleetOrder,
) (models.Fleet, error) {
	universe, err := u.universeRepo.GetByPlanetId(ctx, planet)
	if err != nil {
		return models.Fleet{}, err
	}

	if !universe.ValidCoordinates(order.Destination) {
		return models.Fleet{}, domainerrors.ErrCoordinatesOutOfBound
	}

	speed, err := domainservices.DetermineFleetSpeed(order, universe.Ships)
	if err != nil {
		return models.Fleet{}, err
	}

	moment := u.clock.Now(ctx)

	mutator := generateFleetCreationMutator(order, speed, universe.OccupancyMap)
	result, err := domainservices.AdvanceSimulationUntil(ctx, u.simulator, planet, moment, mutator)
	if err != nil {
		return models.Fleet{}, err
	}

	if result.Fleet == nil {
		return models.Fleet{}, domainerrors.ErrNotFound
	}

	return *result.Fleet, nil

}

func generateFleetCreationMutator(
	order models.FleetOrder,
	speed int,
	occupancy models.OccupancyMap,
) models.SimulationStepResolver {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		flight := models.FleetFlight{
			StartTime:   step.Timestamp,
			Source:      step.Planet.Coordinate,
			Destination: order.Destination,
			Target:      nil,
			Speed:       speed,
		}

		if target, ok := occupancy.UsedSlots[order.Destination]; ok {
			flight.Target = &target
		}

		fleet, err := step.Planet.CreateFleet(order, flight)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		return models.SimulationStepResult{
			Resolution: models.FleetCreated,
			Planet:     step.Planet,
			Fleet:      &fleet,
		}, nil
	}
}
