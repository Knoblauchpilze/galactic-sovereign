package drivenadapters

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	metalMineId = uuid.MustParse("d176e82d-f2ca-4611-996b-c4804096caef")
)

func TestIT_BuildingActionRepository_Create(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	t.Run("creates an action", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		action := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			CurrentLevel: 2,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Version:      9,
			Costs:        []models.BuildingActionCost{},
			Storages:     []models.BuildingActionResourceStorage{},
			Productions:  []models.BuildingActionResourceProduction{},
		}

		err := repo.Create(context.Background(), action)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, action.Id)

		actual, err := repo.Get(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("creates an action with costs", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		action := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			CurrentLevel: 2,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Version:      9,
			Costs: []models.BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   36,
				},
				{
					Resource: crystalResourceId,
					Amount:   798,
				},
			},
			Storages:    []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{},
		}

		err := repo.Create(context.Background(), action)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, action.Id)

		actual, err := repo.Get(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("creates an action with storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		action := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			CurrentLevel: 2,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Version:      9,
			Costs:        []models.BuildingActionCost{},
			Storages: []models.BuildingActionResourceStorage{
				{
					Resource: metalResourceId,
					Storage:  321417,
				},
				{
					Resource: crystalResourceId,
					Storage:  65478,
				},
			},
			Productions: []models.BuildingActionResourceProduction{},
		}

		err := repo.Create(context.Background(), action)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, action.Id)

		actual, err := repo.Get(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("creates an action with productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		action := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			CurrentLevel: 2,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Version:      9,
			Costs:        []models.BuildingActionCost{},
			Storages:     []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{
				{
					Resource:   metalResourceId,
					Production: 147,
				},
				{
					Resource:   crystalResourceId,
					Production: 3254,
				},
			},
		}

		err := repo.Create(context.Background(), action)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, action.Id)

		actual, err := repo.Get(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("returns error when action with same id already exists", func(t *testing.T) {
		_, planet := insertTestBuildingAction(t, conn)

		newAction := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			CurrentLevel: 4,
			DesiredLevel: 5,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Minute),
			Version:      14,
		}

		err := repo.Create(context.Background(), newAction)

		actual, ok := db.AsDatabaseError(err)
		require.True(t, ok, "Actual err: %v", err)
		assert.Equal(t, db.ErrUniqueConstraintViolation, actual.Code, "Actual err: %v", err)
		assertBuildingActionDoesNotExist(t, conn, newAction.Id)
	})
}

func TestIT_BuildingActionRepository_Get(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	t.Run("gets an action", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn)

		actual, err := repo.Get(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("gets an action with costs", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionCost)

		actual, err := repo.Get(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("gets an action with storages", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		actual, err := repo.Get(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("gets an action with productions", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		actual, err := repo.Get(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("returns error when action does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(context.Background(), id)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func TestIT_BuildingActionRepository_ListForPlanet(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	t.Run("lists an action for a planet", func(t *testing.T) {
		action1, planet1 := insertTestBuildingAction(t, conn)
		action2, _ := insertTestBuildingAction(t, conn)

		actual, err := repo.ListForPlanet(context.Background(), planet1.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.GreaterOrEqual(t, len(actual), 1)
		assert.Contains(t, actual, action1)
		assert.NotContains(t, actual, action2)
	})

	t.Run("lists an action with costs for a planet", func(t *testing.T) {
		action1, planet1 := insertTestBuildingAction(t, conn, addBuildingActionCost)
		action2, _ := insertTestBuildingAction(t, conn)

		actual, err := repo.ListForPlanet(context.Background(), planet1.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.GreaterOrEqual(t, len(actual), 1)
		assert.Contains(t, actual, action1)
		assert.NotContains(t, actual, action2)
	})

	t.Run("lists an action with storages for a planet", func(t *testing.T) {
		action1, planet1 := insertTestBuildingAction(t, conn, addBuildingActionStorage)
		action2, _ := insertTestBuildingAction(t, conn)

		actual, err := repo.ListForPlanet(context.Background(), planet1.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.GreaterOrEqual(t, len(actual), 1)
		assert.Contains(t, actual, action1)
		assert.NotContains(t, actual, action2)
	})

	t.Run("lists an action with productions for a planet", func(t *testing.T) {
		action1, planet1 := insertTestBuildingAction(t, conn, addBuildingActionProduction)
		action2, _ := insertTestBuildingAction(t, conn)

		actual, err := repo.ListForPlanet(context.Background(), planet1.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.GreaterOrEqual(t, len(actual), 1)
		assert.Contains(t, actual, action1)
		assert.NotContains(t, actual, action2)
	})
}

func TestIT_BuildingActionRepository_ListBeforeCompletionTime(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	t.Run("returns action when completion is before requested time", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn)

		until := action.CompletedAt.Add(1 * time.Second)

		actual, err := repo.ListBeforeCompletionTime(context.Background(), planet.Id, until)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, 1, len(actual))
		assert.Contains(t, actual, action)
	})

	t.Run("does not return action when completion is after requested time", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn)

		until := action.CompletedAt.Add(-1 * time.Second)

		actual, err := repo.ListBeforeCompletionTime(context.Background(), planet.Id, until)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, 0, len(actual))
	})
}

func TestIT_BuildingActionRepository_Delete(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	t.Run("deletes an action", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn)

		err := repo.Delete(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with costs", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionCost)

		err := repo.Delete(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with storages", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		err := repo.Delete(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with productions", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		err := repo.Delete(context.Background(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, action.Id)
	})

	t.Run("succeeds when the action does not exist", func(t *testing.T) {
		nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

		err := repo.Delete(context.Background(), nonExistingId)
		require.NoError(t, err, "Actual err: %v", err)
	})
}

func TestIT_BuildingActionRepository_DeleteForPlanet(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	t.Run("deletes an action", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn)

		err := repo.DeleteForPlanet(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with costs", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionCost)

		err := repo.DeleteForPlanet(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes action with storages", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		err := repo.DeleteForPlanet(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes action with productions", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		err := repo.DeleteForPlanet(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, planet.Id)
	})
}

func TestIT_BuildingActionRepository_DeleteForPlayer(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	t.Run("deletes an action", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with costs", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionCost)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes action with storages", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes action with productions", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, planet.Id)
	})
}

func TestIT_BuildingActionRepository_CreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	type testCase struct {
		name   string
		action models.BuildingAction
	}

	testCases := []testCase{
		{
			action: models.BuildingAction{
				Id:           uuid.New(),
				Building:     metalMineId,
				CurrentLevel: 26,
				DesiredLevel: 27,
				CreatedAt:    time.Date(2024, 12, 7, 20, 26, 47, 0, time.UTC),
				CompletedAt:  time.Date(2024, 12, 7, 21, 26, 47, 0, time.UTC),
				Version:      17,
				Costs:        []models.BuildingActionCost{},
				Storages:     []models.BuildingActionResourceStorage{},
				Productions:  []models.BuildingActionResourceProduction{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			planet, _, _ := insertTestPlanetForPlayer(t, conn)

			tc.action.Planet = planet.Id

			func() {
				err := repo.Create(context.Background(), tc.action)
				require.NoError(t, err, "Actual err: %v", err)
			}()

			func() {
				actionFromDb, err := repo.Get(context.Background(), tc.action.Id)
				require.NoError(t, err, "Actual err: %v", err)

				assert.Equal(t, tc.action, actionFromDb)
			}()

			func() {
				err := repo.Delete(context.Background(), tc.action.Id)
				require.NoError(t, err, "Actual err: %v", err)
			}()

			assertBuildingActionDoesNotExist(t, conn, tc.action.Id)
		})
	}
}

func newTestBuildingActionRepository(t *testing.T) (drivenports.ForManagingBuildingActions, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewBuildingActionRepository(conn), conn
}

func insertTestBuildingAction(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.BuildingAction),
) (models.BuildingAction, models.Planet) {
	t.Helper()

	planet, _, _ := insertTestPlanetForPlayer(t, conn)

	action := models.BuildingAction{
		Id:           uuid.New(),
		Planet:       planet.Id,
		Building:     metalMineId,
		CurrentLevel: 4,
		DesiredLevel: 5,
		CreatedAt:    someTime,
		CompletedAt:  someTime.Add(1*time.Hour + 2*time.Minute),
		Version:      3,
		// This is intentional: the details (e.g. costs, productions, etc.) are returned as empty
		// slices by the adapter
		Costs:       []models.BuildingActionCost{},
		Storages:    []models.BuildingActionResourceStorage{},
		Productions: []models.BuildingActionResourceProduction{},
	}

	sqlQuery := `INSERT INTO building_action
		(id, planet, building, current_level, desired_level, created_at, completed_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		action.Id,
		action.Planet,
		action.Building,
		action.CurrentLevel,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
		action.Version,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &action)
	}

	return action, planet
}

func addBuildingActionCost(t *testing.T, conn db.Connection, a *models.BuildingAction) {
	t.Helper()

	cost := models.BuildingActionCost{
		Resource: metalResourceId,
		Amount:   rand.Intn(4589),
	}

	sqlQuery := `INSERT INTO building_action_cost (action, resource, amount)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		a.Id,
		cost.Resource,
		cost.Amount,
	)
	require.NoError(t, err, "Actual err: %v", err)

	a.Costs = append(a.Costs, cost)
}

func addBuildingActionStorage(t *testing.T, conn db.Connection, a *models.BuildingAction) {
	t.Helper()

	storage := models.BuildingActionResourceStorage{
		Resource: crystalResourceId,
		Storage:  rand.Intn(65114),
	}

	sqlQuery := `INSERT INTO building_action_resource_storage (action, resource, storage)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		a.Id,
		storage.Resource,
		storage.Storage,
	)
	require.NoError(t, err, "Actual err: %v", err)

	a.Storages = append(a.Storages, storage)
}

func addBuildingActionProduction(t *testing.T, conn db.Connection, a *models.BuildingAction) {
	t.Helper()

	production := models.BuildingActionResourceProduction{
		Resource:   crystalResourceId,
		Production: rand.Intn(7451),
	}

	sqlQuery := `INSERT INTO building_action_resource_production (action, resource, production)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		a.Id,
		production.Resource,
		production.Production,
	)
	require.NoError(t, err, "Actual err: %v", err)

	a.Productions = append(a.Productions, production)
}

func assertBuildingActionExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, 1, value)
}

func assertBuildingActionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionCostDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_cost WHERE action = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionStorageDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_storage WHERE action = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionProductionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_production WHERE action = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
