package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_OccupancyMap_PickPosition(t *testing.T) {
	t.Run("picks a random position", func(t *testing.T) {
		m := OccupancyMap{
			Topology: UniverseTopology{
				Galaxies:     1,
				SolarSystems: 2,
				Orbits:       10,
			},
		}

		c1 := m.PickPosition(uuid.New())
		c2 := m.PickPosition(uuid.New())

		assert.NotEqual(t, c1, c2)
	})

	t.Run("does not pick used coordinate", func(t *testing.T) {
		m := OccupancyMap{
			Topology: UniverseTopology{
				Galaxies:     1,
				SolarSystems: 1,
				Orbits:       2,
			},
			UsedSlots: map[Coordinate]uuid.UUID{
				{Galaxy: 0, SolarSystem: 0, Position: 0}: uuid.New(),
			},
		}

		actual := m.PickPosition(uuid.New())

		expected := Coordinate{Galaxy: 0, SolarSystem: 0, Position: 1}
		assert.Equal(t, expected, actual)
	})

	t.Run("attaches used slot to planet", func(t *testing.T) {
		m := OccupancyMap{
			Topology: UniverseTopology{
				Galaxies:     1,
				SolarSystems: 1,
				Orbits:       1,
			},
		}

		planet := uuid.New()
		m.PickPosition(planet)

		expected := map[Coordinate]uuid.UUID{
			{Galaxy: 0, SolarSystem: 0, Position: 0}: planet,
		}
		assert.Equal(t, expected, m.UsedSlots)
	})
}
