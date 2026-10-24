package domainservices

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

// ResolveSimulationStepThen generates a function closure which takes a simulation step
// and returns a result.
// Upon receiving the step, the closure is checking if it contains a fleet event: if it
// is the case, the event is resolved first and the mutation is not applied. The result
// will be either FleetEventResolved or FleetDeleted.
// In case there are no fleet event, the planet attached to the step is advanced to the
// current time and the complete step is forwarded to the mutation function.
// The mutation is expected to return a simulation step result which can potentially
// include a mutated planet, and a fleet. The resolution should indicate which it is.
// If the mutation returns FleetEventResolved or FleetDeleted, the closure will return
// an error as those are reserved values.
func ResolveSimulationStepThen(
	mutation func(*models.SimulationStep) (models.SimulationStepResult, error),
) func(*models.SimulationStep) (models.SimulationStepResult, error) {
	return func(step *models.SimulationStep) (models.SimulationStepResult, error) {
		if step.NextFleetEvent != nil {
			return resolveFleet(*step.NextFleetEvent)
		}

		err := AdvancePlanetToTime(step.Planet, step.Timestamp)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		result, err := mutation(step)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		switch result.Resolution {
		case models.FleetEventResolved, models.FleetDeleted:
			return models.SimulationStepResult{}, domainerrors.ErrInvalidSimulationStepResult
		}

		return result, nil
	}
}

func resolveFleet(event models.FleetEventContext) (models.SimulationStepResult, error) {
	err := event.Fleet.UpdateToTime(event.Timestamp)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	if event.Planet != nil {
		err = event.Planet.UpdateToTime(event.Timestamp)
		if err != nil {
			return models.SimulationStepResult{}, err
		}
	}

	switch event.Fleet.Schedule.Status {
	case models.StatusOutbound:
		err = event.Fleet.ApplyArrival(event.Planet)
	case models.StatusReturning:
		err = event.Fleet.ApplyReturn(event.Planet)
	default:
		return models.SimulationStepResult{}, domainerrors.ErrUnsupportedFleetStatus
	}
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	resolution := models.FleetEventResolved
	if event.Fleet.Completed() {
		resolution = models.FleetDeleted
	}

	result := models.SimulationStepResult{
		Resolution: resolution,
		Planet:     event.Planet,
		Fleet:      nil,
	}

	if resolution != models.FleetDeleted {
		result.Fleet = &event.Fleet
	}

	return result, err
}
