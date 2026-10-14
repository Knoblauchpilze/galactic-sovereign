package domainservices

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_ResolveSimulationStep(t *testing.T) {
	t.Run("advances planet to time when no fleet event is available", func(t *testing.T) {
		step := models.SimulationStep{
			Timestamp: t3,
			Planet: models.Planet{
				Id:        uuid.New(),
				CreatedAt: t1,
				UpdatedAt: t2,
				Version:   1,
			},
			NextFleetEvent: nil,
		}

		actual, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet: &models.Planet{
				Id:        step.Planet.Id,
				CreatedAt: t1,
				UpdatedAt: t3,
				Version:   2,
			},
			Fleet: nil,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("advances fleet to time when fleet arrives", func(t *testing.T) {
		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t3,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t2,
				Fleet:     fleet,
				Planet:    &target,
			},
		}

		result, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, models.FleetEventResolved, result.Resolution)
		require.NotNil(t, result.Fleet)
		assert.Equal(t, fleet.Id, result.Fleet.Id)
		assert.Equal(t, t2, result.Fleet.UpdatedAt)
		assert.Equal(t, 3, result.Fleet.Version)
	})

	t.Run("advances planet to time when fleet arrives", func(t *testing.T) {
		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t3,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t2,
				Fleet:     fleet,
				Planet:    &target,
			},
		}

		result, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, models.FleetEventResolved, result.Resolution)
		require.NotNil(t, result.Planet)
		assert.Equal(t, target.Id, result.Planet.Id)
		assert.Equal(t, t2, result.Planet.UpdatedAt)
		assert.Equal(t, 2, result.Planet.Version)
	})

	t.Run("applies arrival to fleet when fleet arrives", func(t *testing.T) {
		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t3,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t2,
				Fleet:     fleet,
				Planet:    &target,
			},
		}

		result, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, models.FleetEventResolved, result.Resolution)
		require.NotNil(t, result.Fleet)
		assert.Equal(t, fleet.Id, result.Fleet.Id)
		assert.Equal(t, models.StatusReturning, result.Fleet.Schedule.Status)
	})

	t.Run("applies arrival to planet when fleet arrives", func(t *testing.T) {
		// TODO: No op for now
	})

	t.Run("advances fleet to time when fleet returns", func(t *testing.T) {
		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t4,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t4,
				Fleet:     fleet,
				Planet:    &source,
			},
		}

		result, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, result.Fleet)
		assert.Equal(t, fleet.Id, result.Fleet.Id)
		assert.Equal(t, models.StatusReturning, result.Fleet.Schedule.Status)
		assert.Equal(t, t4, result.Fleet.UpdatedAt)
		assert.Equal(t, 2, result.Fleet.Version)
	})

	t.Run("advances planet to time when fleet returns", func(t *testing.T) {
		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t4,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t4,
				Fleet:     fleet,
				Planet:    &source,
			},
		}

		result, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, result.Planet)
		assert.Equal(t, source.Id, result.Planet.Id)
		assert.Equal(t, t4, result.Planet.UpdatedAt)
		assert.Equal(t, 3, result.Planet.Version)
	})

	t.Run("applies return to fleet when fleet returns", func(t *testing.T) {
		// TODO: No op for now
	})

	t.Run("applies return to planet when fleet returns", func(t *testing.T) {
		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
			Ships: []models.PlanetShip{
				{Ship: lightFighterId, Count: 3},
				{Ship: smallCargoId, Count: 4},
			},
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Ships: []models.FleetShip{
				{Ship: smallCargoId, Count: 1},
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t4,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t4,
				Fleet:     fleet,
				Planet:    &source,
			},
		}

		result, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, result.Planet)
		assert.Equal(t, source.Id, result.Planet.Id)
		expected := []models.PlanetShip{
			{Ship: lightFighterId, Count: 3},
			{Ship: smallCargoId, Count: 5},
		}
		assert.Equal(t, expected, result.Planet.Ships)
	})

	t.Run("marks fleet as completed when fleet returns", func(t *testing.T) {
		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t4,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t4,
				Fleet:     fleet,
				Planet:    &source,
			},
		}

		result, err := ResolveSimulationStep(step)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, models.FleetDeleted, result.Resolution)
		require.NotNil(t, result.Fleet)
		assert.Equal(t, fleet.Id, result.Fleet.Id)
		assert.True(t, result.Fleet.Completed())
	})

	t.Run("returns error when planet already progressed past fleet arrival time", func(t *testing.T) {

		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t3,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusOutbound,
			},
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t4,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t2,
				Fleet:     fleet,
				Planet:    &target,
			},
		}

		_, err := ResolveSimulationStep(step)

		assert.ErrorIs(t, err, domainerrors.ErrPlanetNotUpToDate, "Actual err: %v", err)
	})

	t.Run("returns error when planet already progressed past fleet return time", func(t *testing.T) {
		source := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t5,
			Version:   1,
		}
		target := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t4,
			Version:   1,
		}
		fleet := models.Fleet{
			Id:     uuid.New(),
			Source: source.Id,
			Destination: models.FleetDestination{
				Target: &target.Id,
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2,
				ReturnAt:  t4,
				Status:    models.StatusReturning,
			},
			CreatedAt: t1,
			UpdatedAt: t4,
			Version:   1,
		}

		step := models.SimulationStep{
			Timestamp: t4,
			Planet:    source,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: t4,
				Fleet:     fleet,
				Planet:    &source,
			},
		}

		_, err := ResolveSimulationStep(step)

		assert.ErrorIs(t, err, domainerrors.ErrPlanetNotUpToDate, "Actual err: %v", err)
	})

}
