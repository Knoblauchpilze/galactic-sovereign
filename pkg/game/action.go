package game

import (
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

func ValidateBuildingAction(action persistence.BuildingAction, resources []persistence.PlanetResource, costs []persistence.BuildingCost, buildings []persistence.PlanetBuilding) error {
	if err := validateActionLevel(action); err != nil {
		return err
	}

	if err := validateActionBuilding(action, buildings); err != nil {
		return err
	}

	if err := validateActionCost(resources, costs); err != nil {
		return err
	}

	return nil
}

func validateActionLevel(action persistence.BuildingAction) error {
	if action.CurrentLevel+1 != action.DesiredLevel {
		return errors.NewCode(InvalidActionData)
	}
	return nil
}

func validateActionBuilding(action persistence.BuildingAction, buildings []persistence.PlanetBuilding) error {
	for _, building := range buildings {
		if building.Building == action.Building {
			if building.Level != action.CurrentLevel {
				return errors.NewCode(InvalidBuildingLevel)
			} else {
				return nil
			}
		}
	}

	return errors.NewCode(InvalidBuildingLevel)
}

func validateActionCost(resources []persistence.PlanetResource, costs []persistence.BuildingCost) error {
	temp := make(map[uuid.UUID]persistence.PlanetResource)
	for _, resource := range resources {
		temp[resource.Resource] = resource
	}

	for _, cost := range costs {
		actual, ok := temp[cost.Resource]
		if !ok || actual.Amount < float64(cost.Cost) {
			return errors.NewCode(NotEnoughResources)
		}
	}

	return nil
}
