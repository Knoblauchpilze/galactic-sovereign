package game

import (
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

func ConsolidateBuildingAction(action persistence.BuildingAction, buildings []persistence.PlanetBuilding) persistence.BuildingAction {
	for _, building := range buildings {
		if building.Building == action.Building {
			action.CurrentLevel = building.Level
			break
		}
	}

	action.DesiredLevel = action.CurrentLevel + 1

	return action
}

func ValidateBuildingAction(action persistence.BuildingAction, resources []persistence.PlanetResource, buildings []persistence.PlanetBuilding, costs []persistence.BuildingCost) error {
	if err := validateActionBuilding(action, buildings); err != nil {
		return err
	}

	return validateActionCost(resources, costs)
}

func validateActionBuilding(action persistence.BuildingAction, buildings []persistence.PlanetBuilding) error {
	for _, building := range buildings {
		if action.Building == building.Building {
			return nil
		}
	}

	return errors.NewCode(NoSuchBuilding)
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
