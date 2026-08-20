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
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
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
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})

	t.Run("assigns planet when multiple planets already exist", func(t *testing.T) {
		homeworldId := uuid.New()
		p := Player{
			Id:        uuid.New(),
			Homeworld: homeworldId,
			Planets: []PlayerPlanet{
				{
					Id: homeworldId,
					Coordinate: Coordinate{Galaxy: 2,
						SolarSystem: 12,
						Position:    14,
					},
				},
			},
		}

		actual := p.Colonize(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.False(t, actual.Homeworld)
		assert.Equal(t, "colony", actual.Name)
		assert.NotEqual(t, actual.Id, p.Homeworld)

		assert.Equal(t, []uuid.UUID{homeworldId, actual.Id}, p.Planets)
	})
}

func assertPlanetIsWithinRange(
	t *testing.T,
	actual PlayerPlanet,
	topology UniverseTopology,
) {
	t.Helper()

	assert.Greater(t, actual.Coordinate.Galaxy, 0)
	assert.Less(t, actual.Coordinate.Galaxy, topology.Galaxies)
	assert.Greater(t, actual.Coordinate.SolarSystem, 0)
	assert.Less(t, actual.Coordinate.SolarSystem, topology.SolarSystems)
	assert.Greater(t, actual.Coordinate.Position, 0)
	assert.Less(t, actual.Coordinate.Position, topology.Orbits)
}
