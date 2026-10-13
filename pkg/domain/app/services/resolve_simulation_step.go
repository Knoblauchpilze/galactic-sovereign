package domainservices

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

func ResolveSimulationStep(step models.SimulationStep) (models.SimulationStepResult, error) {
	if step.NextFleetEvent != nil {
		return resolveFleet(*step.NextFleetEvent)
	}

	err := AdvancePlanetToTime(&step.Planet, step.Timestamp)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	result := models.SimulationStepResult{
		Planet: &step.Planet,
		Fleet:  nil,
	}

	return result, nil
}

func resolveFleet(event models.FleetEventContext) (models.SimulationStepResult, error) {
	err := event.Fleet.UpdateToTime(event.Timestamp)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	err = event.Planet.UpdateToTime(event.Timestamp)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	switch event.Fleet.Schedule.Status {
	case models.StatusOutbound:
		err = event.Fleet.ApplyArrival(event.Planet)
	case models.StatusReturning:
		err = event.Fleet.ApplyReturn(event.Planet)
	default:
		return models.SimulationStepResult{}, domainerrors.ErrUnsupportedFleetStatus
	}

	result := models.SimulationStepResult{
		Planet: event.Planet,
		Fleet:  &event.Fleet,
	}

	return result, err
}
