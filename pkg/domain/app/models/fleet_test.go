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

	t.Run("returns error when returning fleet has already been updated past the arrival time", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t4,
			Version:   1,
		}

		err := f.UpdateToTime(t5)

		assert.ErrorIs(t, err, domainerrors.ErrFleetAlreadyArrived, "Actual err is: %v", err)
	})

	t.Run("does not mark fleet as returning when arrival time is not reached", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t3,
				ReturnAt:  t4,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		err := f.UpdateToTime(t2)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, StatusOutbound, f.Schedule.Status)
	})

	t.Run("marks fleet as returning when fleet arrives at its destination", func(t *testing.T) {
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

		err := f.UpdateToTime(t2)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, StatusReturning, f.Schedule.Status)
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

func TestUnit_Fleet_ApplyArrival(t *testing.T) {

	t.Run("returns error when destination does not match target for undirected fleet", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     nil,
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		p := Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		err := f.ApplyArrival(&p)

		assert.ErrorIs(t, err, domainerrors.ErrWrongFleetTarget, "Actual err is: %v", err)
	})

	t.Run("returns error when destination does not match target for directed fleet", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     new(uuid.New()),
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		p := Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		err := f.ApplyArrival(&p)

		assert.ErrorIs(t, err, domainerrors.ErrWrongFleetTarget, "Actual err is: %v", err)
	})

	t.Run("returns error when destination is not provided for directed fleet", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     new(uuid.New()),
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		err := f.ApplyArrival(nil)

		assert.ErrorIs(t, err, domainerrors.ErrWrongFleetTarget, "Actual err is: %v", err)
	})

	t.Run("returns error when fleet is not up to date with arrival time", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     nil,
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		err := f.ApplyArrival(nil)

		assert.ErrorIs(t, err, domainerrors.ErrFleetNotUpToDate, "Actual err is: %v", err)
	})

	t.Run("returns error when planet is not up to date with arrival time", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     new(uuid.New()),
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		p := Planet{
			Id:        *f.Destination.Target,
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		err := f.ApplyArrival(&p)

		assert.ErrorIs(t, err, domainerrors.ErrPlanetNotUpToDate, "Actual err is: %v", err)
	})

	t.Run("bumps target updated at field", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     new(uuid.New()),
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		p := Planet{
			Id:        *f.Destination.Target,
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		err := f.ApplyArrival(&p)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t2, f.UpdatedAt)
		assert.Equal(t, 2, p.Version)
	})
}

func TestUnit_Fleet_ApplyReturn(t *testing.T) {
	t.Run("returns error when destination does not match target", func(t *testing.T) {
		f := Fleet{
			Id:     uuid.New(),
			Source: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     nil,
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		p := Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		err := f.ApplyReturn(&p)

		assert.ErrorIs(t, err, domainerrors.ErrWrongFleetTarget, "Actual err is: %v", err)
	})

	t.Run("returns error when destination is nil", func(t *testing.T) {
		f := Fleet{
			Id:     uuid.New(),
			Source: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     nil,
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		err := f.ApplyReturn(nil)

		assert.ErrorIs(t, err, domainerrors.ErrWrongFleetTarget, "Actual err is: %v", err)
	})

	t.Run("returns error when fleet is not up to date with return time", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     nil,
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		p := Planet{
			Id:        f.Source,
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}

		err := f.ApplyReturn(&p)

		assert.ErrorIs(t, err, domainerrors.ErrFleetNotUpToDate, "Actual err is: %v", err)
	})

	t.Run("returns error when planet is not up to date with return time", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     new(uuid.New()),
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}

		p := Planet{
			Id:        f.Source,
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		err := f.ApplyReturn(&p)

		assert.ErrorIs(t, err, domainerrors.ErrPlanetNotUpToDate, "Actual err is: %v", err)
	})

	t.Run("registers fleet ships as planet ships on return", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     new(uuid.New()),
			},
			Ships: []FleetShip{
				{Ship: lightFighterId, Count: 3},
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}

		p := Planet{
			Id: f.Source,
			Ships: []PlanetShip{
				{Ship: lightFighterId, Count: 2},
				{Ship: smallCargoId, Count: 1},
			},
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}

		err := f.ApplyReturn(&p)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetShip{
			{Ship: lightFighterId, Count: 5},
			{Ship: smallCargoId, Count: 1},
		}
		assert.Equal(t, expected, p.Ships)
	})

	t.Run("bumps target updated at field", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Destination: FleetDestination{
				Coordinate: Coordinate{},
				Target:     new(uuid.New()),
			},
			Schedule: FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t3,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}

		p := Planet{
			Id:        f.Source,
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}

		err := f.ApplyReturn(&p)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t3, f.UpdatedAt)
		assert.Equal(t, 2, p.Version)
	})
}

func TestUnit_Fleet_Completed(t *testing.T) {
	t.Run("returns false when fleet was just created", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t3,
				ReturnAt:  t5,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
		}

		assert.False(t, f.Completed())
	})

	t.Run("returns false when fleet is on its outbound journey", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t3,
				ReturnAt:  t5,
				Status:    StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
		}

		assert.False(t, f.Completed())
	})

	t.Run("returns false when fleet completed its outbound journey", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t3,
				ReturnAt:  t5,
				Status:    StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t3,
		}

		assert.False(t, f.Completed())
	})

	t.Run("returns false when fleet is on its returning journey", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t3,
				ReturnAt:  t5,
				Status:    StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t4,
		}

		assert.False(t, f.Completed())
	})

	t.Run("returns true when fleet completed its return journey", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t3,
				ReturnAt:  t5,
				Status:    StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t5,
		}

		assert.True(t, f.Completed())
	})

	t.Run("returns true when fleet is past its return time", func(t *testing.T) {
		f := Fleet{
			Id: uuid.New(),
			Schedule: FleetFlightSchedule{
				ArrivalAt: t3,
				ReturnAt:  t5,
				Status:    StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t5.Add(1 * time.Minute),
		}

		assert.True(t, f.Completed())
	})
}
