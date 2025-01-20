package game

import (
	"math"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
)

const resourceUnitsPerHour float64 = 2500.0

func findResourceByName(name string, resources []persistence.Resource) *persistence.Resource {
	for _, resource := range resources {
		if resource.Name == name {
			return &resource
		}
	}

	return nil
}

func findResourceRequirementByName(name string, resources []persistence.Resource, costs []persistence.BuildingActionCost) (float64, error) {
	resource := findResourceByName(name, resources)
	if resource == nil {
		return 0.0, errors.NewCode(NoSuchResource)
	}

	for _, cost := range costs {
		if cost.Resource == resource.Id {
			return float64(cost.Amount), nil
		}
	}

	return 0.0, nil
}

func buildingCompletionTimeFromCost(resources []persistence.Resource, costs []persistence.BuildingActionCost) (time.Duration, error) {
	// https://ogame.fandom.com/wiki/Buildings
	metal, err := findResourceRequirementByName("metal", resources, costs)
	if err != nil {
		return 0, err
	}
	crystal, err := findResourceRequirementByName("crystal", resources, costs)
	if err != nil {
		return 0, err
	}

	buildTimeHour := (metal + crystal) / resourceUnitsPerHour

	nanoSeconds := math.Ceil(buildTimeHour * float64(time.Hour.Nanoseconds()))

	return time.Duration(nanoSeconds), nil
}
