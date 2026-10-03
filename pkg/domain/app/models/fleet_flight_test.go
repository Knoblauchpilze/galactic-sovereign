package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnit_FleetFlight_Duration(t *testing.T) {
	t.Run("calculates duration for flight between two different galaxies", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      1,
				SolarSystem: 120,
				Position:    2,
			},
			Destination: Coordinate{
				Galaxy:      2,
				SolarSystem: 56,
				Position:    14,
			},
			Speed: 2000,
		}

		actual := f.Duration()

		expected := 35010 * time.Second
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates duration for flight between two different solar systems", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      2,
				SolarSystem: 10,
				Position:    0,
			},
			Destination: Coordinate{
				Galaxy:      2,
				SolarSystem: 198,
				Position:    13,
			},
			Speed: 6500,
		}

		actual := f.Duration()

		// Non rounded value: 19694.43352187 seconds
		expected := 19694 * time.Second
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates duration for flight within same solar system", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      0,
				SolarSystem: 23,
				Position:    9,
			},
			Destination: Coordinate{
				Galaxy:      0,
				SolarSystem: 23,
				Position:    6,
			},
			Speed: 15000,
		}

		actual := f.Duration()

		// Non rounded value: 2889.091291826 seconds
		expected := 2889 * time.Second
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates duration for flight to same coordinates", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      6,
				SolarSystem: 451,
				Position:    5,
			},
			Destination: Coordinate{
				Galaxy:      6,
				SolarSystem: 451,
				Position:    5,
			},
			Speed: 13000,
		}

		actual := f.Duration()

		// Non rounded value: 227.060785531 seconds
		expected := 227 * time.Second
		assert.Equal(t, expected, actual)
	})
}

func TestUnit_FleetFlight_ArrivalTime(t *testing.T) {
	t.Run("calculates arrival time for flight between two different galaxies", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      1,
				SolarSystem: 120,
				Position:    2,
			},
			Destination: Coordinate{
				Galaxy:      2,
				SolarSystem: 56,
				Position:    14,
			},
			Speed: 2000,
		}

		actual := f.ArrivalTime()

		expected := someTime.Add(35010 * time.Second)
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates arrival time for flight between two different solar systems", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      2,
				SolarSystem: 10,
				Position:    0,
			},
			Destination: Coordinate{
				Galaxy:      2,
				SolarSystem: 198,
				Position:    13,
			},
			Speed: 6500,
		}

		actual := f.ArrivalTime()

		// Non rounded value: 19694.43352187 seconds
		expected := someTime.Add(19694 * time.Second)
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates arrival time for flight within same solar system", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      0,
				SolarSystem: 23,
				Position:    9,
			},
			Destination: Coordinate{
				Galaxy:      0,
				SolarSystem: 23,
				Position:    6,
			},
			Speed: 15000,
		}

		actual := f.ArrivalTime()

		// Non rounded value: 2889.091291826 seconds
		expected := someTime.Add(2889 * time.Second)
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates arrival time for flight to same coordinates", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      6,
				SolarSystem: 451,
				Position:    5,
			},
			Destination: Coordinate{
				Galaxy:      6,
				SolarSystem: 451,
				Position:    5,
			},
			Speed: 13000,
		}

		actual := f.ArrivalTime()

		// Non rounded value: 227.060785531 seconds
		expected := someTime.Add(227 * time.Second)
		assert.Equal(t, expected, actual)
	})
}

func TestUnit_FleetFlight_ReturnTime(t *testing.T) {
	t.Run("calculates return time for flight between two different galaxies", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      1,
				SolarSystem: 120,
				Position:    2,
			},
			Destination: Coordinate{
				Galaxy:      2,
				SolarSystem: 56,
				Position:    14,
			},
			Speed: 2000,
		}

		actual := f.ReturnTime()

		expected := someTime.Add(2 * 35010 * time.Second)
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates return time for flight between two different solar systems", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      2,
				SolarSystem: 10,
				Position:    0,
			},
			Destination: Coordinate{
				Galaxy:      2,
				SolarSystem: 198,
				Position:    13,
			},
			Speed: 6500,
		}

		actual := f.ReturnTime()

		// Non rounded value: 19694.43352187 seconds
		expected := someTime.Add(2 * 19694 * time.Second)
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates return time for flight within same solar system", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      0,
				SolarSystem: 23,
				Position:    9,
			},
			Destination: Coordinate{
				Galaxy:      0,
				SolarSystem: 23,
				Position:    6,
			},
			Speed: 15000,
		}

		actual := f.ReturnTime()

		// Non rounded value: 2889.091291826 seconds
		expected := someTime.Add(2 * 2889 * time.Second)
		assert.Equal(t, expected, actual)
	})

	t.Run("calculates return time for flight to same coordinates", func(t *testing.T) {
		f := FleetFlight{
			StartTime: someTime,
			Source: Coordinate{
				Galaxy:      6,
				SolarSystem: 451,
				Position:    5,
			},
			Destination: Coordinate{
				Galaxy:      6,
				SolarSystem: 451,
				Position:    5,
			},
			Speed: 13000,
		}

		actual := f.ReturnTime()

		// Non rounded value: 227.060785531 seconds
		expected := someTime.Add(2 * 227 * time.Second)
		assert.Equal(t, expected, actual)
	})
}
