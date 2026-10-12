package domainservices

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
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

	result := models.SimulationStepResult{
		Planet: nil,
		Fleet:  &event.Fleet,
	}

	return result, nil
}
