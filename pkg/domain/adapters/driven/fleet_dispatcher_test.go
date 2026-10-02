package drivenadapters

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	t2 = someOtherTime
	t3 = t2.Add(1 * time.Hour)
	t4 = t2.Add(2 * time.Hour)
)

func TestIT_FleetDispatcher_GetBehavior(t *testing.T) {
	testCases := []struct {
		name      string
		generator func(t *testing.T, conn database.Connection) models.Planet
	}{
		{
			name: "planet",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn)
				return planet
			},
		},
		{
			name: "planet with resources",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
				return planet
			},
		},
		{
			name: "planet with ships",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
				return planet
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adapter, conn := newTestFleetDispatcher(t)

			planet := tc.generator(t, conn)
			fleet := generateSampleFleet(planet)

			var captured models.Planet
			mutator := func(p *models.Planet) (models.Fleet, error) {
				captured = *p
				p.Version++
				return fleet, nil
			}

			actual, err := adapter.Dispatch(t.Context(), planet.Id, mutator)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, planet, captured)
			assert.Equal(t, fleet, actual)
		})
	}
}

func TestIT_FleetDispatcher_MutateBehavior(t *testing.T) {
	t.Run("returns dispatched fleet", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		expected := generateSampleFleet(planet)

		creator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			p.UpdatedAt = yetAnotherTime
			p.Version++
			return expected
		})

		actual, err := adapter.Dispatch(t.Context(), planet.Id, creator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("persists dispatched fleet", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)
		require.NotEqual(t, 326, planet.Fields)

		fleet := generateSampleFleet(planet)

		creator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			p.Fields = 326
			p.UpdatedAt = yetAnotherTime
			p.Version++

			return fleet
		})

		_, err := adapter.Dispatch(t.Context(), planet.Id, creator)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, fleet, actual)
	})

	t.Run("persists mutated planet", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)
		require.NotEqual(t, 326, planet.Fields)

		creator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			p.Fields = 326
			p.UpdatedAt = yetAnotherTime
			p.Version++
			return generateSampleFleet(*p)
		})

		_, err := adapter.Dispatch(t.Context(), planet.Id, creator)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		expected := planet
		expected.Fields = 326
		expected.UpdatedAt = yetAnotherTime
		expected.Version++
		assert.Equal(t, expected, actual)
	})

	t.Run("persists mutated planet resources", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 5874, planet.Resources[0].Amount)

		fleet := generateSampleFleet(planet)

		mutator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			p.Resources[0].Amount = 5874
			p.Version++
			return fleet
		})

		returned, err := adapter.Dispatch(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, fleet, returned)
		assertPlanetResourceAmount(t, conn, planet.Id, crystalResourceId, 5874)
	})

	t.Run("does not delete existing planet resource", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		resource := planet.Resources[0].Resource
		amount := planet.Resources[0].Amount

		fleet := generateSampleFleet(planet)

		mutator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			p.Resources = []models.PlanetResource{}
			p.Version++
			return fleet
		})

		returned, err := adapter.Dispatch(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, fleet, returned)
		assertPlanetResourceAmount(t, conn, planet.Id, resource, amount)
	})

	t.Run("persists mutated planet ships", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		require.NotEqual(t, 776, planet.Ships[0].Count)

		creator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			p.Ships[0].Count = 776
			p.Version++
			return generateSampleFleet(*p)
		})

		_, err := adapter.Dispatch(t.Context(), planet.Id, creator)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetShipCount(t, conn, planet.Id, lightFighterId, 776)
	})

	t.Run("does not delete existing planet ship", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		ship := planet.Ships[0].Ship
		count := planet.Ships[0].Count

		creator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			p.Ships = []models.PlanetShip{}
			p.Version++
			return generateSampleFleet(*p)
		})

		_, err := adapter.Dispatch(t.Context(), planet.Id, creator)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetShipCount(t, conn, planet.Id, ship, count)
	})

	t.Run("returns error when planet does not exist", func(t *testing.T) {
		adapter, _ := newTestFleetDispatcher(t)

		creator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			p.UpdatedAt = yetAnotherTime
			p.Version++
			return generateSampleFleet(*p)
		})

		_, err := adapter.Dispatch(t.Context(), uuid.New(), creator)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when mutator does not update version", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		fleet := generateSampleFleet(planet)

		creator := generateFleetCreator(func(p *models.Planet) models.Fleet {
			return fleet
		})

		_, err := adapter.Dispatch(t.Context(), planet.Id, creator)

		assert.ErrorIs(t, err, domainerrors.ErrMutationWithoutVersionBump, "Actual err: %v", err)
		assertFleetDoesNotExist(t, conn, planet.Id)
	})
}

func TestIT_FleetDispatcher_Dispatch_Concurrency(t *testing.T) {
	t.Run("blocks concurrent fleet creation for same planet", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 9876.0, planet.Resources[0].Amount)

		fleetA := generateSampleFleet(planet)
		fleetB := generateSampleFleet(planet)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		enteredB := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			_, err := adapter.Dispatch(
				t.Context(),
				planet.Id,
				func(p *models.Planet) (models.Fleet, error) {
					close(enteredA)
					<-releaseA
					p.UpdatedAt = t3
					p.Version++
					return fleetA, nil
				},
			)
			doneA <- err
		}()

		<-enteredA

		blockingCtx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		go func() {
			_, err := adapter.Dispatch(
				blockingCtx,
				planet.Id,
				func(p *models.Planet) (models.Fleet, error) {
					close(enteredB)
					p.UpdatedAt = t4
					p.Version++
					return fleetB, nil
				})
			doneB <- err
		}()

		err := <-doneB
		require.ErrorIs(
			t, err, context.DeadlineExceeded, "Expected deadline exceeded, got err: %v", err,
		)

		select {
		case <-enteredB:
			t.Fatalf("second mutation reached callback even though first mutation held lock")
		default:
		}

		close(releaseA)

		err = <-doneA
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadFleetFromDb(t, conn, fleetA.Id)
		assert.Equal(t, fleetA, actual)
		assertFleetDoesNotExist(t, conn, fleetB.Id)
		actualPlanet := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, t3, actualPlanet.UpdatedAt)
	})

	t.Run("does not block concurrent fleet creation for different planets", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, yetAnotherTime, planetA.UpdatedAt)
		planetB, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 8765.0, planetB.Resources[0].Amount)

		fleetA := generateSampleFleet(planetA)
		fleetB := generateSampleFleet(planetB)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			_, err := adapter.Dispatch(
				t.Context(),
				planetA.Id,
				func(p *models.Planet) (models.Fleet, error) {
					close(enteredA)
					<-releaseA
					p.UpdatedAt = yetAnotherTime
					p.Version++
					return fleetA, nil
				})
			doneA <- err
		}()

		<-enteredA

		go func() {
			_, err := adapter.Dispatch(
				t.Context(),
				planetB.Id,
				func(p *models.Planet) (models.Fleet, error) {
					p.Resources[0].Amount = 8765.0
					p.Version++
					return fleetB, nil
				})
			doneB <- err
		}()

		select {
		case err := <-doneB:
			require.NoError(t, err, "Actual err: %v", err)
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("mutation on second planet was blocked while first planet was locked")
		}

		assertPlanetResourceAmount(t, conn, planetB.Id, crystalResourceId, 8765.0)
		assertFleetExists(t, conn, fleetB.Id)

		close(releaseA)

		err := <-doneA
		require.NoError(t, err, "Actual err: %v", err)

		assertFleetExists(t, conn, fleetA.Id)
	})

	t.Run("waiting mutation respects context timeout", func(t *testing.T) {
		adapter, conn := newTestFleetDispatcher(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		fleetA := generateSampleFleet(planet)
		fleetB := generateSampleFleet(planet)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		enteredB := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			_, err := adapter.Dispatch(
				t.Context(),
				planet.Id,
				func(p *models.Planet) (models.Fleet, error) {
					close(enteredA)
					<-releaseA
					p.Resources[0].Amount = 3333.0
					p.Version++
					return fleetA, nil
				})
			doneA <- err
		}()

		<-enteredA

		blockingCtx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		go func() {
			_, err := adapter.Dispatch(
				blockingCtx,
				planet.Id,
				func(p *models.Planet) (models.Fleet, error) {
					close(enteredB)
					p.Resources[0].Amount = 2222.0
					p.Version++
					return fleetB, nil
				})
			doneB <- err
		}()

		errB := <-doneB
		require.ErrorIs(
			t, errB, context.DeadlineExceeded, "Expected deadline exceeded, got err: %v", errB,
		)

		select {
		case <-enteredB:
			t.Fatalf("waiting mutation callback should not be reached before lock is released")
		default:
		}

		close(releaseA)

		errA := <-doneA
		require.NoError(t, errA, "Actual err: %v", errA)

		actual := loadFleetFromDb(t, conn, fleetA.Id)
		assert.Equal(t, fleetA, actual)
		assertFleetDoesNotExist(t, conn, fleetB.Id)
		actualPlanet := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, 3333.0, actualPlanet.Resources[0].Amount)
	})
}

func newTestFleetDispatcher(t *testing.T) (*FleetDispatcher, database.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewFleetDispatcher(conn), conn
}

func generateFleetCreator(modifier func(p *models.Planet) models.Fleet) drivenports.FleetCreator {
	return func(p *models.Planet) (models.Fleet, error) {
		f := modifier(p)
		return f, nil
	}
}

func generateSampleFleet(p models.Planet) models.Fleet {
	return models.Fleet{
		Id:     uuid.New(),
		Player: p.Player,
		Source: p.Id,
		Destination: models.Coordinate{
			Galaxy:      p.Coordinate.Galaxy,
			SolarSystem: p.Coordinate.SolarSystem + 1,
			Position:    p.Coordinate.Position,
		},
		Ships: []models.FleetShip{
			{
				Ship:  lightFighterId,
				Count: 12,
			},
			{
				Ship:  smallCargoId,
				Count: 3,
			},
		},
		CreatedAt: p.UpdatedAt,
		ArrivalAt: p.UpdatedAt.Add(1 * time.Hour),
		ReturnAt:  p.UpdatedAt.Add(2 * time.Hour),
		UpdatedAt: p.UpdatedAt,
		Version:   0,
	}
}

func loadFleetFromDb(t *testing.T, conn database.Connection, id uuid.UUID) models.Fleet {
	t.Helper()

	tx, err := conn.BeginTx(t.Context())
	require.NoError(t, err, "Actual err: %v", err)
	defer tx.Close(t.Context())

	fleet, err := loadFleetAndDetails(t.Context(), tx, id)
	require.NoError(t, err, "Actual err: %v", err)

	return fleet
}

func assertFleetExists(t *testing.T, conn database.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT id FROM fleet WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, id, value)
}

func assertFleetDoesNotExist(t *testing.T, conn database.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM fleet WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
