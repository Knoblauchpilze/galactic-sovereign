package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Coordinate_DistanceTo(t *testing.T) {
	t.Run("calculates distance when coordinates are in a different galaxy", func(t *testing.T) {
		lhs := Coordinate{Galaxy: 2, SolarSystem: 235, Position: 12}
		rhs := Coordinate{Galaxy: 5, SolarSystem: 111, Position: 2}

		assert.Equal(t, 60_000, lhs.DistanceTo(rhs))
		assert.Equal(t, 60_000, rhs.DistanceTo(lhs))

		lhs = Coordinate{Galaxy: 7, SolarSystem: 24, Position: 3}
		rhs = Coordinate{Galaxy: 2, SolarSystem: 120, Position: 3}

		assert.Equal(t, 100_000, lhs.DistanceTo(rhs))
		assert.Equal(t, 100_000, rhs.DistanceTo(lhs))
	})

	t.Run("calculates distance when coordinates are in a different solar system", func(t *testing.T) {
		lhs := Coordinate{Galaxy: 0, SolarSystem: 29, Position: 2}
		rhs := Coordinate{Galaxy: 0, SolarSystem: 34, Position: 17}

		assert.Equal(t, 3175, lhs.DistanceTo(rhs))
		assert.Equal(t, 3175, rhs.DistanceTo(lhs))

		lhs = Coordinate{Galaxy: 0, SolarSystem: 354, Position: 5}
		rhs = Coordinate{Galaxy: 0, SolarSystem: 25, Position: 0}

		assert.Equal(t, 33_955, lhs.DistanceTo(rhs))
		assert.Equal(t, 33_955, rhs.DistanceTo(lhs))
	})

	t.Run("calculates distance for coordinates in the same solar system", func(t *testing.T) {
		lhs := Coordinate{Galaxy: 3, SolarSystem: 412, Position: 2}
		rhs := Coordinate{Galaxy: 3, SolarSystem: 412, Position: 5}

		assert.Equal(t, 1015, lhs.DistanceTo(rhs))
		assert.Equal(t, 1015, rhs.DistanceTo(lhs))

		lhs = Coordinate{Galaxy: 8, SolarSystem: 12, Position: 0}
		rhs = Coordinate{Galaxy: 8, SolarSystem: 12, Position: 14}

		assert.Equal(t, 1070, lhs.DistanceTo(rhs))
		assert.Equal(t, 1070, rhs.DistanceTo(lhs))
	})

	t.Run("returns minimum distance for same coordinate", func(t *testing.T) {
		lhs := Coordinate{Galaxy: 3, SolarSystem: 412, Position: 2}
		assert.Equal(t, 5, lhs.DistanceTo(lhs))

		lhs = Coordinate{Galaxy: 8, SolarSystem: 12, Position: 0}
		assert.Equal(t, 5, lhs.DistanceTo(lhs))
	})
}
