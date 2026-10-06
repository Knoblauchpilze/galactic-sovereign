package drivenadapters

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlanetDeleter_GetBehavior(t *testing.T) {
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
			name: "planet with resource productions",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
				return planet
			},
		},
		{
			name: "planet with resource productions for building",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
				return planet
			},
		},
		{
			name: "planet with resource storages",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
				return planet
			},
		},
		{
			name: "planet with buildings",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
				return planet
			},
		},
		{
			name: "planet with buildings with speedup",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuildingWithShipSpeedup)
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
		{
			name: "planet with building action",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuildingAction)
				return planet
			},
		},
		{
			name: "planet with ship actions",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShipAction)
				return planet
			},
		},
		{
			name: "planet with multiple ship actions are sorted",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn)
				action1 := insertTestShipActionForPlanet(t, conn, planet.Id, func(t *testing.T, a *models.ShipAction) {
					a.CreatedAt = planet.CreatedAt.Add(2 * time.Hour)
				})
				action2 := insertTestShipActionForPlanet(t, conn, planet.Id, func(t *testing.T, a *models.ShipAction) {
					a.CreatedAt = planet.CreatedAt.Add(1 * time.Hour)
				})

				planet.ShipActions = []models.ShipAction{action2, action1}

				return planet
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adapter, conn := newTestPlanetDeleter(t)

			planet := tc.generator(t, conn)

			var captured models.Planet
			deleter := func(p *models.Planet) error {
				captured = *p
				return nil
			}

			err := adapter.Delete(t.Context(), planet.Id, deleter)
			require.NoError(t, err, "Actual err: %v", err)

			assert.Equal(t, planet, captured)
		})
	}
}

func TestIT_PlanetDeleter_Delete(t *testing.T) {
	t.Run("deletes planet", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		player, _ := insertTestPlayerInUniverse(t, conn)

		err := adapter.Delete(t.Context(), player.Homeworld, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, player.Homeworld)
		assertPlanetIsNotHomeworld(t, conn, player.Homeworld)
	})

	t.Run("deletes planet with resources", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with storages", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with productions", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with production for building", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with buildings", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetBuildingDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with buildings with ship speedup", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuildingWithShipSpeedup)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetBuildingDoesNotExist(t, conn, planet.Id)
		assertBuildingShipSpeedupValue(t, conn, planet.Buildings[0].Building, *planet.Buildings[0].ShipSpeedup)
	})

	t.Run("deletes planet with ships", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetShipDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with building action", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuildingAction)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		require.NotNil(t, planet.BuildingAction)
		assertBuildingActionDoesNotExist(t, conn, planet.BuildingAction.Id)
	})

	t.Run("deletes planet with ship action", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShipAction)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		err := adapter.Delete(t.Context(), planet.Id, dummyDelete)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertShipActionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("does not delete planet and returns error when planet is source for a fleet", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet1, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet1, addFleetShip)

		err := adapter.Delete(t.Context(), planet1.Id, dummyDelete)

		assert.ErrorIs(t, err, domainerrors.ErrFleetInFlight, "Actual err: %v", err)
		assertPlanetExists(t, conn, planet1.Id)
		assertFleetExists(t, conn, fleet.Id)
	})

	t.Run("does not delete planet and returns error when planet has a fleet targeting it", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet1, _, _ := insertTestPlanetForPlayer(t, conn)
		planet2, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet1, addFleetShip)
		fleet.Destination = models.FleetDestination{
			Coordinate: planet2.Coordinate,
			Target:     &planet2.Id,
		}
		upsertFleetDestination(t, conn, fleet)

		err := adapter.Delete(t.Context(), planet2.Id, dummyDelete)

		assert.ErrorIs(t, err, domainerrors.ErrFleetInFlight, "Actual err: %v", err)
		assertPlanetExists(t, conn, planet1.Id)
		assertPlanetExists(t, conn, planet2.Id)
		assertFleetExists(t, conn, fleet.Id)
	})
}

func TestIT_PlanetDeleter_Delete_Concurrency(t *testing.T) {
	t.Run("blocks concurrent deletion for same planet", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 9876.0, planet.Resources[0].Amount)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		enteredB := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			err := adapter.Delete(t.Context(), planet.Id, func(p *models.Planet) error {
				close(enteredA)
				<-releaseA
				p.UpdatedAt = yetAnotherTime
				p.Version++
				return nil
			})
			doneA <- err
		}()

		<-enteredA

		blockingCtx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		go func() {
			err := adapter.Delete(blockingCtx, planet.Id, func(p *models.Planet) error {
				close(enteredB)
				p.Resources[0].Amount = 9877.0
				p.Version++
				return nil
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

		assertPlanetDoesNotExist(t, conn, planet.Id)
	})

	t.Run("does not block concurrent mutation for different planets", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, yetAnotherTime, planetA.UpdatedAt)
		planetB, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 8765.0, planetB.Resources[0].Amount)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			err := adapter.Delete(t.Context(), planetA.Id, func(p *models.Planet) error {
				close(enteredA)
				<-releaseA
				return nil
			})
			doneA <- err
		}()

		<-enteredA

		go func() {
			err := adapter.Delete(t.Context(), planetB.Id, func(p *models.Planet) error {
				return nil
			})
			doneB <- err
		}()

		select {
		case err := <-doneB:
			require.NoError(t, err, "Actual err: %v", err)
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("mutation on second planet was blocked while first planet was locked")
		}

		assertPlanetDoesNotExist(t, conn, planetB.Id)
		assertPlanetExists(t, conn, planetA.Id)

		close(releaseA)

		err := <-doneA
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planetA.Id)
	})

	t.Run("waiting mutation respects context timeout", func(t *testing.T) {
		adapter, conn := newTestPlanetDeleter(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		enteredB := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			err := adapter.Delete(t.Context(), planet.Id, func(p *models.Planet) error {
				close(enteredA)
				<-releaseA
				return nil
			})
			doneA <- err
		}()

		<-enteredA

		blockingCtx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		go func() {
			err := adapter.Delete(blockingCtx, planet.Id, func(p *models.Planet) error {
				close(enteredB)
				return nil
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

		assertPlanetExists(t, conn, planet.Id)

		close(releaseA)

		errA := <-doneA
		require.NoError(t, errA, "Actual err: %v", errA)

		assertPlanetDoesNotExist(t, conn, planet.Id)
	})
}

func newTestPlanetDeleter(t *testing.T) (*PlanetDeleter, database.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewPlanetDeleter(conn), conn
}

func dummyDelete(*models.Planet) error {
	return nil
}
