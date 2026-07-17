package models

import (
	"testing"

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

		c1 := m.PickPosition()
		c2 := m.PickPosition()

		assert.NotEqual(t, c1, c2)
	})

	t.Run("does not pick used coordinate", func(t *testing.T) {
		m := OccupancyMap{
			Topology: UniverseTopology{
				Galaxies:     1,
				SolarSystems: 1,
				Orbits:       2,
			},
			UsedSlots: map[Coordinate]struct{}{
				Coordinate{Galaxy: 0, SolarSystem: 0, Position: 0}: struct{}{},
			},
		}

		actual := m.PickPosition()

		expected := Coordinate{Galaxy: 0, SolarSystem: 0, Position: 1}
		assert.Equal(t, expected, actual)
	})
}
