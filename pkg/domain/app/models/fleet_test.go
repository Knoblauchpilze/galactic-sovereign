package models

import (
	"testing"
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	t1 = time.Date(2026, time.October, 3, 20, 31, 52, 0, time.UTC)
)

func TestUnit_Fleet_New(t *testing.T) {
	t.Run("creates a fleet with colonize mission", func(t *testing.T) {
		origin := FleetOrigin{
			Source:     uuid.New(),
			Coordinate: Coordinate{1, 3, 5},
			Player:     uuid.New(),
		}
		order := FleetOrder{
			Mission:     MissionColonize,
			Destination: Coordinate{1, 3, 6},
			Ships: []FleetShip{
				{Ship: uuid.New(), Count: 14},
			},
		}
		flight := FleetFlight{
			StartTime:   t1,
			Source:      origin.Coordinate,
			Destination: order.Destination,
			Target:      nil,
			Speed:       3100,
		}

		fleet, err := NewFleet(origin, order, flight)
		require.NoError(t, err, "Actual err: %v", err)

		// Non rounded value: 6311.881439492 seconds
		flightDuration := 6311 * time.Second
		expected := Fleet{
			Id:      fleet.Id,
			Player:  origin.Player,
			Source:  origin.Source,
			Mission: order.Mission,
			Destination: FleetDestination{
				Coordinate: flight.Destination,
				Target:     nil,
			},
			Ships:     order.Ships,
			CreatedAt: t1,
			ArrivalAt: t1.Add(flightDuration),
			ReturnAt:  t1.Add(2 * flightDuration),
			UpdatedAt: t1,
			Version:   0,
		}
		assert.Equal(t, expected, fleet)
	})

	t.Run("returns error when colonize mission has a target", func(t *testing.T) {
		origin := FleetOrigin{
			Source:     uuid.New(),
			Coordinate: Coordinate{1, 3, 5},
			Player:     uuid.New(),
		}
		order := FleetOrder{
			Mission:     MissionColonize,
			Destination: Coordinate{1, 3, 6},
			Ships: []FleetShip{
				{Ship: uuid.New(), Count: 14},
			},
		}
		flight := FleetFlight{
			StartTime:   t1,
			Source:      origin.Coordinate,
			Destination: order.Destination,
			Target:      new(uuid.New()),
			Speed:       3100,
		}

		_, err := NewFleet(origin, order, flight)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidFleetConfiguration, "Actual err: %v", err)
	})

	t.Run("creates a fleet with transport mission", func(t *testing.T) {
		origin := FleetOrigin{
			Source:     uuid.New(),
			Coordinate: Coordinate{1, 3, 5},
			Player:     uuid.New(),
		}
		order := FleetOrder{
			Mission:     MissionTransport,
			Destination: Coordinate{1, 3, 6},
			Ships: []FleetShip{
				{Ship: uuid.New(), Count: 14},
			},
		}
		flight := FleetFlight{
			StartTime:   t1,
			Source:      origin.Coordinate,
			Destination: order.Destination,
			Target:      new(uuid.New()),
			Speed:       3100,
		}

		fleet, err := NewFleet(origin, order, flight)
		require.NoError(t, err, "Actual err: %v", err)

		// Non rounded value: 6311.881439492 seconds
		flightDuration := 6311 * time.Second
		expected := Fleet{
			Id:      fleet.Id,
			Player:  origin.Player,
			Source:  origin.Source,
			Mission: order.Mission,
			Destination: FleetDestination{
				Coordinate: flight.Destination,
				Target:     flight.Target,
			},
			Ships:     order.Ships,
			CreatedAt: t1,
			ArrivalAt: t1.Add(flightDuration),
			ReturnAt:  t1.Add(2 * flightDuration),
			UpdatedAt: t1,
			Version:   0,
		}
		assert.Equal(t, expected, fleet)
	})

	t.Run("returns error when transport mission does not have a target", func(t *testing.T) {
		origin := FleetOrigin{
			Source:     uuid.New(),
			Coordinate: Coordinate{1, 3, 5},
			Player:     uuid.New(),
		}
		order := FleetOrder{
			Mission:     MissionTransport,
			Destination: Coordinate{1, 3, 6},
			Ships: []FleetShip{
				{Ship: uuid.New(), Count: 14},
			},
		}
		flight := FleetFlight{
			StartTime:   t1,
			Source:      origin.Coordinate,
			Destination: order.Destination,
			Target:      nil,
			Speed:       3100,
		}

		_, err := NewFleet(origin, order, flight)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidFleetConfiguration, "Actual err: %v", err)
	})

	t.Run("return error when fleet has no ship", func(t *testing.T) {
		origin := FleetOrigin{
			Source:     uuid.New(),
			Coordinate: Coordinate{1, 3, 5},
			Player:     uuid.New(),
		}
		order := FleetOrder{
			Mission:     MissionTransport,
			Destination: Coordinate{1, 3, 6},
			Ships:       []FleetShip{},
		}
		flight := FleetFlight{
			StartTime:   t1,
			Source:      origin.Coordinate,
			Destination: order.Destination,
			Target:      new(uuid.New()),
			Speed:       3100,
		}

		_, err := NewFleet(origin, order, flight)

		assert.ErrorIs(t, err, domainerrors.ErrNoShipInFleet, "Actual err: %v", err)
	})
}
