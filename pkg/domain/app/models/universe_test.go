package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	metalResourceId   = uuid.MustParse("b4419b6b-b3bf-4576-aa92-055283addbc8")
	crystalResourceId = uuid.MustParse("cd2ac9aa-9968-4ff5-b746-88f1f810fbb3")
)

func TestUnit_Universe_CreatePlanet(t *testing.T) {
	playerId := uuid.New()

	t.Run("creates a homeworld belonging to the player", func(t *testing.T) {
		u := sampleUniverse()

		beforeCreation := time.Now()
		actual := u.CreatePlanet(playerId, true)

		assert.Equal(t, playerId, actual.Player)
		assert.Equal(t, "homeworld", actual.Name)
		assert.True(t, actual.Homeworld)
		assert.True(t, beforeCreation.Before(actual.CreatedAt))
		assert.Equal(t, actual.CreatedAt, actual.UpdatedAt)
		assert.Zero(t, actual.Version)
		assert.Nil(t, actual.BuildingAction)
	})

	t.Run("creates a colony belonging to the player", func(t *testing.T) {
		u := sampleUniverse()

		beforeCreation := time.Now()
		actual := u.CreatePlanet(playerId, false)

		assert.Equal(t, playerId, actual.Player)
		assert.Equal(t, "colony", actual.Name)
		assert.False(t, actual.Homeworld)
		assert.True(t, beforeCreation.Before(actual.CreatedAt))
		assert.Equal(t, actual.CreatedAt, actual.UpdatedAt)
		assert.Zero(t, actual.Version)
		assert.Nil(t, actual.BuildingAction)
	})

	t.Run("assigns start amount for each resource", func(t *testing.T) {
		u := sampleUniverse()

		actual := u.CreatePlanet(playerId, false)

		expected := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   145,
			},
			{
				Resource: crystalResourceId,
				Amount:   325,
			},
		}
		assert.Equal(t, expected, actual.Resources)
	})

	t.Run("assigns start storage for each resource", func(t *testing.T) {
		u := sampleUniverse()

		actual := u.CreatePlanet(playerId, false)

		expected := []PlanetResourceStorage{
			{
				Resource: metalResourceId,
				Storage:  123,
			},
			{
				Resource: crystalResourceId,
				Storage:  421,
			},
		}
		assert.Equal(t, expected, actual.Storages)
	})

	t.Run("assigns start production for each resource", func(t *testing.T) {
		u := sampleUniverse()

		actual := u.CreatePlanet(playerId, false)

		expected := []PlanetResourceProduction{
			{
				Resource:   metalResourceId,
				Production: 985,
			},
			{
				Resource:   crystalResourceId,
				Production: 752,
			},
		}
		assert.Equal(t, expected, actual.Productions)
	})

	t.Run("creates each building with level 0", func(t *testing.T) {
		u := sampleUniverse()

		actual := u.CreatePlanet(playerId, false)

		expected := []PlanetBuilding{
			{
				Building: u.Buildings[0].Id,
				Level:    0,
			},
			{
				Building: u.Buildings[1].Id,
				Level:    0,
			},
		}
		assert.Equal(t, expected, actual.Buildings)
	})

	t.Run("creates each ship with level 0", func(t *testing.T) {
		u := sampleUniverse()

		actual := u.CreatePlanet(playerId, false)

		expected := []PlanetShip{
			{
				Ship:  u.Ships[0].Id,
				Count: 0,
			},
			{
				Ship:  u.Ships[1].Id,
				Count: 0,
			},
		}
		assert.Equal(t, expected, actual.Ships)
	})

	t.Run("creates a planet at a free spot", func(t *testing.T) {
		u := sampleUniverse()
		u.OccupancyMap = OccupancyMap{
			Topology: UniverseTopology{
				Galaxies:     1,
				SolarSystems: 1,
				Orbits:       2,
			},
			UsedSlots: map[Coordinate]struct{}{
				{Galaxy: 0, SolarSystem: 0, Position: 1}: {},
			},
		}

		actual := u.CreatePlanet(playerId, false)

		expected := Coordinate{
			Galaxy:      0,
			SolarSystem: 0,
			Position:    0,
		}
		assert.Equal(t, expected, actual.Coordinate)
	})
}

func sampleResources() []Resource {
	return []Resource{
		{
			Id:              metalResourceId,
			StartAmount:     145,
			StartStorage:    123,
			StartProduction: 985,
		},
		{
			Id:              crystalResourceId,
			StartAmount:     325,
			StartStorage:    421,
			StartProduction: 752,
		},
	}
}

func sampleBuildings() []Building {
	return []Building{
		{
			Id: uuid.New(),
			Costs: []BuildingCost{
				{
					Resource: metalResourceId,
					Cost:     550,
					Progress: 26.4,
				},
			},
		},
		{
			Id: uuid.New(),
			Storages: []BuildingResourceStorage{
				{
					Resource: metalResourceId,
					Scale:    1500.0,
					Progress: 26.4,
				},
			},
		},
	}
}

func sampleShips() []Ship {
	return []Ship{
		{
			Id: uuid.New(),
			Costs: []ShipCost{
				{
					Resource: metalResourceId,
					Cost:     987,
				},
			},
		},
		{
			Id: uuid.New(),
			Costs: []ShipCost{
				{
					Resource: crystalResourceId,
					Cost:     1234,
				},
			},
		},
	}
}

func sampleOccupancyMap() OccupancyMap {
	return OccupancyMap{
		Topology: UniverseTopology{
			Galaxies:     5,
			SolarSystems: 281,
			Orbits:       26,
		},
		UsedSlots: make(map[Coordinate]struct{}),
	}
}

func sampleUniverse() Universe {
	occupancyMap := sampleOccupancyMap()
	return Universe{
		Id:           uuid.New(),
		Topology:     occupancyMap.Topology,
		Resources:    sampleResources(),
		Buildings:    sampleBuildings(),
		Ships:        sampleShips(),
		OccupancyMap: occupancyMap,
	}
}
