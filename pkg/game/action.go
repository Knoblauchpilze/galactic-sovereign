package game

import (
	"math"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

func DetermineBuildingActionCost(
	action persistence.BuildingAction,
	baseCosts []persistence.BuildingCost,
) []persistence.BuildingActionCost {
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

func DetermineBuildingActionResourceProduction(
	action persistence.BuildingAction,
	baseProductions []persistence.BuildingResourceProduction,
) []persistence.BuildingActionResourceProduction {
	var productions []persistence.BuildingActionResourceProduction

	// https://ogame.fandom.com/wiki/Metal_Mine#Production
	// https://ogame.fandom.com/wiki/Crystal_Mine#Production
	levelAsFloat := float64(action.DesiredLevel)

	for _, baseProduction := range baseProductions {
		resourceProduction := math.Floor(float64(baseProduction.Base) * levelAsFloat * math.Pow(baseProduction.Progress, levelAsFloat))

		production := persistence.BuildingActionResourceProduction{
			Action:     action.Id,
			Resource:   baseProduction.Resource,
			Production: int(resourceProduction),
		}
		productions = append(productions, production)
	}

	return productions
}

func DetermineBuildingActionResourceStorage(
	action persistence.BuildingAction,
	baseStorages []persistence.BuildingResourceStorage,
) []persistence.BuildingActionResourceStorage {
	var storages []persistence.BuildingActionResourceStorage

	// https://ogame.fandom.com/wiki/Metal_Storage
	// https://ogame.fandom.com/wiki/Crystal_Storage
	// https://ogame.fandom.com/wiki/Deuterium_Tank
	levelAsFloat := float64(action.DesiredLevel)

	for _, baseStorage := range baseStorages {
		// The original form was modified from storage = C * e^(B * level)
		// to fit the form storage = C * C1^level.
		resourceStorage := math.Floor(float64(baseStorage.Base) * math.Floor(baseStorage.Scale*math.Pow(baseStorage.Progress, levelAsFloat)))

		storage := persistence.BuildingActionResourceStorage{
			Action:   action.Id,
			Resource: baseStorage.Resource,
			Storage:  int(resourceStorage),
		}
		storages = append(storages, storage)
	}

	return storages
}

func ConsolidateBuildingActionLevel(
	action persistence.BuildingAction,
	buildings []persistence.PlanetBuilding,
) persistence.BuildingAction {
	for _, building := range buildings {
		if building.Building == action.Building {
			action.CurrentLevel = building.Level
			break
		}
	}

	action.DesiredLevel = action.CurrentLevel + 1

	return action
}

func ConsolidateBuildingActionCompletionTime(
	action persistence.BuildingAction,
	resources []persistence.Resource,
	costs []persistence.BuildingActionCost,
) (persistence.BuildingAction, error) {
	completionTime, err := buildingCompletionTimeFromCost(resources, costs)
	action.CompletedAt = action.CreatedAt.Add(completionTime)

	return action, err
}

func ValidateBuildingAction(
	action persistence.BuildingAction,
	resources []persistence.PlanetResource,
	buildings []persistence.PlanetBuilding,
	costs []persistence.BuildingActionCost,
) error {
	if err := validateActionBuilding(action, buildings); err != nil {
		return err
	}

	return validateActionCost(resources, costs)
}

func validateActionBuilding(
	action persistence.BuildingAction,
	buildings []persistence.PlanetBuilding,
) error {
	for _, building := range buildings {
		if action.Building == building.Building {
			return nil
		}
	}

	return errors.NewCode(NoSuchBuilding)
}

func validateActionCost(
	resources []persistence.PlanetResource,
	costs []persistence.BuildingActionCost,
) error {
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
