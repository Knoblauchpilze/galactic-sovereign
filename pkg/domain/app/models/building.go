package models

import (
	"math"
	"time"

	"github.com/google/uuid"
)

// TODO: A new function should create a building action for a certain level
type Building struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time

	Costs       []BuildingCost
	Productions []BuildingResourceProduction
	Storages    []BuildingResourceStorage
}

type BuildingCost struct {
	Resource uuid.UUID
	Cost     int
	Progress float64
}

type BuildingResourceProduction struct {
	Resource uuid.UUID
	Base     int
	Progress float64
}

type BuildingResourceStorage struct {
	Resource uuid.UUID
	Base     int
	Scale    float64
	Progress float64
}

func (b Building) CreateBuildingAction(planet Planet, building uuid.UUID) (BuildingAction, error) {
	pb, err := findBuildingById(planet.Buildings, building)
	if err != nil {
		return BuildingAction{}, err
	}

	desiredLevel := pb.Level + 1

	action := BuildingAction{
		Id:           uuid.New(),
		Planet:       planet.Id,
		Building:     building,
		CurrentLevel: pb.Level,
		DesiredLevel: desiredLevel,

		CreatedAt: time.Now(),
		// TODO: Calculate the completion time
		CompletedAt: time.Time{},

		Version: 0,

		Costs:       b.determineActionCost(desiredLevel),
		Storages:    b.determineActionResourceStorage(desiredLevel),
		Productions: b.determineActionResourceProduction(desiredLevel),
	}
	return action, nil
}

func (b Building) determineActionCost(
	desiredLevel int,
) []BuildingActionCost {
	costs := []BuildingActionCost{}

	for _, baseCost := range b.Costs {
		resourceCost := math.Floor(float64(baseCost.Cost) * math.Pow(baseCost.Progress, float64(desiredLevel-1)))

		cost := BuildingActionCost{
			Resource: baseCost.Resource,
			Amount:   int(resourceCost),
		}
		costs = append(costs, cost)
	}

	return costs
}

func (b Building) determineActionResourceProduction(
	desiredLevel int,
) []BuildingActionResourceProduction {
	productions := []BuildingActionResourceProduction{}

	// https://ogame.fandom.com/wiki/Metal_Mine#Production
	// https://ogame.fandom.com/wiki/Crystal_Mine#Production
	levelAsFloat := float64(desiredLevel)

	for _, baseProduction := range b.Productions {
		resourceProduction := math.Floor(float64(baseProduction.Base) * levelAsFloat * math.Pow(baseProduction.Progress, levelAsFloat))

		production := BuildingActionResourceProduction{
			Resource:   baseProduction.Resource,
			Production: int(resourceProduction),
		}
		productions = append(productions, production)
	}

	return productions
}

func (b Building) determineActionResourceStorage(
	desiredLevel int,
) []BuildingActionResourceStorage {
	storages := []BuildingActionResourceStorage{}

	// https://ogame.fandom.com/wiki/Metal_Storage
	// https://ogame.fandom.com/wiki/Crystal_Storage
	// https://ogame.fandom.com/wiki/Deuterium_Tank
	levelAsFloat := float64(desiredLevel)

	for _, baseStorage := range b.Storages {
		// The original form was modified from storage = C * e^(B * level)
		// to fit the form storage = C * C1^level.
		resourceStorage := math.Floor(float64(baseStorage.Base) * math.Floor(baseStorage.Scale*math.Pow(baseStorage.Progress, levelAsFloat)))

		storage := BuildingActionResourceStorage{
			Resource: baseStorage.Resource,
			Storage:  int(resourceStorage),
		}
		storages = append(storages, storage)
	}

	return storages
}

func findBuildingById(buildings []PlanetBuilding, id uuid.UUID) (PlanetBuilding, error) {
	for _, b := range buildings {
		if b.Building == id {
			return b, nil
		}
	}

	return PlanetBuilding{}, errBuildingNotFound
}
