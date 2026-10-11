package domainservices

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

func ResolveSimulationStep(step models.SimulationStep) (models.SimulationStepResult, error) {
	if step.NextFleetEvent != nil {
		return models.SimulationStepResult{}, errors.ErrNotImplemented
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
