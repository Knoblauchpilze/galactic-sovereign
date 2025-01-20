package game

import (
	"math"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

func DetermineBuildingActionCost(action persistence.BuildingAction, baseCosts []persistence.BuildingCost) []persistence.BuildingActionCost {
	var costs []persistence.BuildingActionCost
	for _, baseCost := range baseCosts {
		resourceCost := math.Floor(float64(baseCost.Cost) * math.Pow(baseCost.Progress, float64(action.DesiredLevel-1)))

		cost := persistence.BuildingActionCost{
			Action:   action.Id,
			Resource: baseCost.Resource,
			Amount:   int(resourceCost),
		}
		costs = append(costs, cost)
	}

	return costs
}

func DetermineBuildingActionResourceProduction(action persistence.BuildingAction, baseProductions []persistence.BuildingResourceProduction) []persistence.BuildingActionResourceProduction {
	var productions []persistence.BuildingActionResourceProduction
	for _, baseProduction := range baseProductions {
		resourceProduction := math.Floor(float64(baseProduction.Base) * math.Pow(baseProduction.Progress, float64(action.DesiredLevel-1)))

		production := persistence.BuildingActionResourceProduction{
			Action:     action.Id,
			Resource:   baseProduction.Resource,
			Production: int(resourceProduction),
		}
		productions = append(productions, production)
	}

	return productions
}

func ConsolidateBuildingActionLevel(action persistence.BuildingAction, buildings []persistence.PlanetBuilding) persistence.BuildingAction {
	for _, building := range buildings {
		if building.Building == action.Building {
			action.CurrentLevel = building.Level
			break
		}
	}

	action.DesiredLevel = action.CurrentLevel + 1

	return action
}

func ConsolidateBuildingActionCompletionTime(action persistence.BuildingAction, resources []persistence.Resource, costs []persistence.BuildingActionCost) (persistence.BuildingAction, error) {
	completionTime, err := buildingCompletionTimeFromCost(resources, costs)
	action.CompletedAt = action.CreatedAt.Add(completionTime)

	return action, err
}

func ValidateBuildingAction(action persistence.BuildingAction, resources []persistence.PlanetResource, buildings []persistence.PlanetBuilding, costs []persistence.BuildingActionCost) error {
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

func validateActionCost(resources []persistence.PlanetResource, costs []persistence.BuildingActionCost) error {
	temp := make(map[uuid.UUID]persistence.PlanetResource)
	for _, resource := range resources {
		temp[resource.Resource] = resource
	}

	for _, cost := range costs {
		actual, ok := temp[cost.Resource]
		if !ok || actual.Amount < float64(cost.Amount) {
			return errors.NewCode(NotEnoughResources)
		}
	}

	return nil
}
