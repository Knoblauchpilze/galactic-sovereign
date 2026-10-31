package drivenadapters

import (
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	t1 = time.Date(2026, time.October, 31, 16, 21, 23, 0, time.UTC)
	t2 = time.Date(2026, time.October, 31, 17, 21, 23, 0, time.UTC)
	t3 = time.Date(2026, time.October, 31, 18, 21, 23, 0, time.UTC)
	t4 = time.Date(2026, time.October, 31, 19, 21, 23, 0, time.UTC)
)

func TestUnit_SimulationStepContext_NewPlanetStepContext(t *testing.T) {
	t.Run("captures planet version when available", func(t *testing.T) {
		planet := models.Planet{
			Id:      uuid.New(),
			Version: 3,
		}

		ctx := newPlanetStepContext(planet, true, dummyResolver)

		assert.Equal(t, 3, ctx.initialPlanetVersion)
		assert.Equal(t, 0, ctx.initialFleetVersion)
	})

	t.Run("configures planet deletion as configured", func(t *testing.T) {
		planet := models.Planet{
			Id:      uuid.New(),
			Version: 3,
		}

		ctx := newPlanetStepContext(planet, false, dummyResolver)
		assert.False(t, ctx.allowPlanetDeletion)

		ctx = newPlanetStepContext(planet, true, dummyResolver)
		assert.True(t, ctx.allowPlanetDeletion)
	})
}

func TestUnit_SimulationStepContext_NewFleetStepContext(t *testing.T) {
	t.Run("captures fleet version when available", func(t *testing.T) {
		event := models.FleetEventContext{
			Timestamp: someTime,
			Fleet: models.Fleet{
				Id:      uuid.New(),
				Version: 4,
			},
			Planet: nil,
		}

		ctx := newFleetStepContext(event, dummyResolver)

		assert.Equal(t, 0, ctx.initialPlanetVersion)
		assert.Equal(t, 4, ctx.initialFleetVersion)
	})

	t.Run("captures planet and fleet version from fleet event when available", func(t *testing.T) {
		event := models.FleetEventContext{
			Timestamp: someTime,
			Fleet: models.Fleet{
				Id:      uuid.New(),
				Version: 4,
			},
			Planet: &models.Planet{
				Id:      uuid.New(),
				Version: 36,
			},
		}

		ctx := newFleetStepContext(event, dummyResolver)

		assert.Equal(t, 36, ctx.initialPlanetVersion)
		assert.Equal(t, 4, ctx.initialFleetVersion)
	})

	t.Run("disables planet deletion by default", func(t *testing.T) {
		event := models.FleetEventContext{
			Timestamp: someTime,
			Fleet: models.Fleet{
				Id:      uuid.New(),
				Version: 4,
			},
			Planet: nil,
		}

		ctx := newFleetStepContext(event, dummyResolver)

		assert.False(t, ctx.allowPlanetDeletion)
	})
}

func TestUnit_SimulationStepContext_Resolve(t *testing.T) {
	t.Run("forwards resolver error to caller", func(t *testing.T) {
		planet := models.Planet{
			Id:      uuid.New(),
			Version: 3,
		}

		errSample := errors.New("stubbed error")
		resolver := func(_ models.SimulationStep) (models.SimulationStepResult, error) {
			return models.SimulationStepResult{}, errSample
		}

		ctx := newPlanetStepContext(planet, true, resolver)
		_, err := ctx.resolve(t.Context(), nil, t1)

		assert.ErrorIs(t, err, errSample, "Actual err: %v", err)
	})

	t.Run("returns error when resolver requests planet deletion but is not allowed", func(t *testing.T) {
		planet := models.Planet{
			Id:      uuid.New(),
			Version: 3,
		}

		result := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}

		resolver, _ := generateCapturingResolver(result)

		ctx := newPlanetStepContext(planet, false, resolver)

		_, err := ctx.resolve(t.Context(), nil, t2)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})
}

func TestIT_SimulationStepContext_Resolve(t *testing.T) {
	t.Run("calls resolver with correct step for planet update", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     new(planet),
			Fleet:      nil,
		}
		result.Planet.Version++

		resolver, captor := generateCapturingResolver(result)

		ctx := newPlanetStepContext(planet, true, resolver)

		actual, err := ctx.resolve(t.Context(), tx, t1)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SimulationStep{
			Timestamp:      t1,
			Planet:         new(planet),
			NextFleetEvent: nil,
		}
		assert.Equal(t, expected, *captor)
		assert.Equal(t, result, actual)
	})

	t.Run("calls resolver with correct step for fleet update", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet)

		timestamp := fleet.Schedule.ArrivalAt
		event := models.FleetEventContext{
			Timestamp: timestamp,
			Fleet:     fleet,
			Planet:    nil,
		}

		result := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Planet:     nil,
			Fleet:      new(fleet),
		}
		result.Fleet.Version++

		resolver, captor := generateCapturingResolver(result)

		ctx := newFleetStepContext(event, resolver)

		actual, err := ctx.resolve(t.Context(), tx, timestamp)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SimulationStep{
			Timestamp:      timestamp,
			Planet:         nil,
			NextFleetEvent: new(event),
		}
		assert.Equal(t, expected, *captor)
		assert.Equal(t, result, actual)
	})

	t.Run("returns error when resolver requests planet deletion but is not allowed", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		result := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}

		resolver, _ := generateCapturingResolver(result)

		ctx := newPlanetStepContext(planet, false, resolver)

		_, err := ctx.resolve(t.Context(), tx, t2)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})

	t.Run("accepts deletion when resolver requests planet deletion", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		result := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}

		resolver, _ := generateCapturingResolver(result)

		ctx := newPlanetStepContext(planet, true, resolver)

		actual, err := ctx.resolve(t.Context(), tx, t2)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, result, actual)
	})
}

func newTestConnectionAndTransaction(t *testing.T) (database.Connection, database.Transaction) {
	t.Helper()

	conn := newTestConnection(t)
	tx, err := conn.BeginTx(t.Context())
	require.NoError(t, err, "Actual err: %v", err)

	t.Cleanup(func() {
		tx.Close(t.Context())
	})

	return conn, tx
}

func generateCapturingResolver(result models.SimulationStepResult) (models.SimulationStepResolver, *models.SimulationStep) {
	var captor models.SimulationStep

	resolver := func(step models.SimulationStep) (models.SimulationStepResult, error) {
		captor = step
		return result, nil
	}

	return resolver, &captor
}

func dummyResolver(_ models.SimulationStep) (models.SimulationStepResult, error) {
	return models.SimulationStepResult{
		Resolution: models.PlanetUpdated,
		Planet:     &models.Planet{},
		Fleet:      nil,
	}, nil
}
