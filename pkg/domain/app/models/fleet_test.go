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
	t2 = time.Date(2026, time.October, 12, 8, 36, 8, 0, time.UTC)
	t3 = time.Date(2026, time.October, 12, 9, 36, 8, 0, time.UTC)
	t4 = time.Date(2026, time.October, 12, 10, 36, 8, 0, time.UTC)
	t5 = time.Date(2026, time.October, 12, 11, 36, 8, 0, time.UTC)
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
			Ships: order.Ships,
			Schedule: FleetFlightSchedule{
				ArrivalAt: t1.Add(flightDuration),
				ReturnAt:  t1.Add(2 * flightDuration),
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
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
			Ships: order.Ships,
			Schedule: FleetFlightSchedule{
				ArrivalAt: t1.Add(flightDuration),
				ReturnAt:  t1.Add(2 * flightDuration),
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
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

func TestUnit_Fleet_UpdateToTime(t *testing.T) {

	t.Run("does not change resource when update time is before already computed time", func(t *testing.T) {
		f := Fleet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}

		err := f.UpdateToTime(t2)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t3, f.UpdatedAt)
		assert.Equal(t, 1, f.Version)
	})

	t.Run("returns error when outbound fleet has arrival time before the update time", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		err := f.UpdateToTime(t3)

		assert.ErrorIs(t, err, domainerrors.ErrFleetNotUpToDate, "Actual err is: %v", err)
	})

	t.Run("returns error when returning fleet has return time before the update time", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}

		err := f.UpdateToTime(t5)

		assert.ErrorIs(t, err, domainerrors.ErrFleetNotUpToDate, "Actual err is: %v", err)
	})

	t.Run("bumps updated at field", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		err := f.UpdateToTime(t3)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t3, f.UpdatedAt)
	})
}
