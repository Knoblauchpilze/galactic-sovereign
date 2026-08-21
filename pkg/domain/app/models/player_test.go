package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_Player_CreateHomeworld(t *testing.T) {
	u := sampleUniverse()

	t.Run("creates a homeworld belonging to the player", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []PlayerPlanet{},
		}

		actual := p.CreateHomeworld(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.True(t, actual.Homeworld)
		assert.Equal(t, "homeworld", actual.Name)
		assert.Equal(t, actual.Id, p.Homeworld)

		assert.Len(t, p.Planets, 1)
		assert.Equal(t, actual.Id, p.Planets[0].Id)
		assert.Equal(t, actual.Name, p.Planets[0].Name)
		assertPlanetIsWithinRange(t, p.Planets[0], u.Topology)
	})

	t.Run("assigns planet when slice is nil", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: nil,
		}

		actual := p.CreateHomeworld(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.True(t, actual.Homeworld)
		assert.Equal(t, "homeworld", actual.Name)
		assert.Equal(t, actual.Id, p.Homeworld)

		assert.Len(t, p.Planets, 1)
		assert.Equal(t, actual.Id, p.Planets[0].Id)
		assert.Equal(t, actual.Name, p.Planets[0].Name)
		assertPlanetIsWithinRange(t, p.Planets[0], u.Topology)
	})
}

func TestUnit_Player_Colonize(t *testing.T) {
	u := sampleUniverse()

	t.Run("creates a planet belonging to the player", func(t *testing.T) {
		p := Player{
			Id:        uuid.New(),
			Homeworld: uuid.New(),
			Planets:   []PlayerPlanet{},
		}

		actual := p.Colonize(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.False(t, actual.Homeworld)
		assert.Equal(t, "colony", actual.Name)
		assert.NotEqual(t, actual.Id, p.Homeworld)

		assert.Len(t, p.Planets, 1)
		assert.Equal(t, actual.Id, p.Planets[0].Id)
		assert.Equal(t, actual.Name, p.Planets[0].Name)
		assertPlanetIsWithinRange(t, p.Planets[0], u.Topology)
	})

	t.Run("assigns planet when slice is nil", func(t *testing.T) {
		p := Player{
			Id:        uuid.New(),
			Homeworld: uuid.New(),
			Planets:   nil,
		}

		actual := p.Colonize(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.False(t, actual.Homeworld)
		assert.Equal(t, "colony", actual.Name)
		assert.NotEqual(t, actual.Id, p.Homeworld)

		assert.Len(t, p.Planets, 1)
		assert.Equal(t, actual.Id, p.Planets[0].Id)
		assert.Equal(t, actual.Name, p.Planets[0].Name)
		assertPlanetIsWithinRange(t, p.Planets[0], u.Topology)
	})

	t.Run("assigns planet when multiple planets already exist", func(t *testing.T) {
		homeworld := PlayerPlanet{
			Id:   uuid.New(),
			Name: "homeworld",
			Coordinate: Coordinate{Galaxy: 2,
				SolarSystem: 12,
				Position:    14,
			},
		}

		p := Player{
			Id:        uuid.New(),
			Homeworld: homeworld.Id,
			Planets:   []PlayerPlanet{homeworld},
		}

		actual := p.Colonize(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.False(t, actual.Homeworld)
		assert.Equal(t, "colony", actual.Name)
		assert.NotEqual(t, actual.Id, p.Homeworld)

		assert.Len(t, p.Planets, 2)
		assert.Equal(t, homeworld, p.Planets[0])
		assert.Equal(t, actual.Id, p.Planets[1].Id)
		assert.Equal(t, actual.Name, p.Planets[1].Name)
		assertPlanetIsWithinRange(t, p.Planets[1], u.Topology)
	})
}

func assertPlanetIsWithinRange(
	t *testing.T,
	actual PlayerPlanet,
	topology UniverseTopology,
) {
	t.Helper()

	assert.GreaterOrEqual(t, actual.Coordinate.Galaxy, 0, "Galaxy is not within bounds: %v, %v", actual.Coordinate, topology)
	assert.Less(t, actual.Coordinate.Galaxy, topology.Galaxies, "Galaxy is not within bounds: %v, %v", actual.Coordinate, topology)
	assert.GreaterOrEqual(t, actual.Coordinate.SolarSystem, 0, "Solar system is not within bounds: %v, %v", actual.Coordinate, topology)
	assert.Less(t, actual.Coordinate.SolarSystem, topology.SolarSystems, "Solar system is not within bounds: %v, %v", actual.Coordinate, topology)
	assert.GreaterOrEqual(t, actual.Coordinate.Position, 0, "Orbit is not within bounds: %v, %v", actual.Coordinate, topology)
	assert.Less(t, actual.Coordinate.Position, topology.Orbits, "Orbit is not within bounds: %v, %v", actual.Coordinate, topology)
}
