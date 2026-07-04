package drivenadapters

import (
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
	planetRepo := NewPlanetRepository(conn)

	t.Run("creates an action", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs:        []models.BuildingActionCost{},
			Storages:     []models.BuildingActionResourceStorage{},
			Productions:  []models.BuildingActionResourceProduction{},
		}
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, *planet.BuildingAction, *actual.BuildingAction)
	})

	t.Run("creates an action with costs", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
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
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, *planet.BuildingAction, *actual.BuildingAction)
	})

	t.Run("creates an action with storages", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
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
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, *planet.BuildingAction, *actual.BuildingAction)
	})

	t.Run("creates an action with productions", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
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
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, *planet.BuildingAction, *actual.BuildingAction)
	})

	t.Run("updates planet version and updated at", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs: []models.BuildingActionCost{
				{
					Resource: crystalResourceId,
					Amount:   1000,
				},
			},
			Storages:    []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{},
		}
		newTime := time.Date(2026, time.June, 20, 14, 01, 17, 0, time.UTC)
		require.NotEqual(t, planet.UpdatedAt, newTime)
		planet.UpdatedAt = newTime
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet.Version, actual.Version)
		assert.Equal(t, newTime, actual.UpdatedAt)
	})

	t.Run("updates planet resources", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 1000.0, planet.Resources[0].Amount)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			// This should probably reflect the value of 1000 set with the planet
			// but as this does not play a role in the test it is left out
			Costs:       []models.BuildingActionCost{},
			Storages:    []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{},
		}
		planet.Version++

		planet.Resources[0].Amount -= 1000

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		expectedAmount := planet.Resources[0].Amount
		assertPlanetResourceAmount(t, conn, planet.Id, crystalResourceId, float64(expectedAmount))
	})

	t.Run("updates planet storages", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		require.NotEqual(t, planet.Storages[0].Storage, 5000)
		planet.Storages[0].Storage = 5000
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs:        []models.BuildingActionCost{},
			// This should probably reflect the storage of 5000 set with the planet
			// but as this does not play a role in the test it is left out
			Storages:    []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{},
		}
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		assertPlanetResourceStorage(t, conn, planet.Id, crystalResourceId, 5000)
	})

	t.Run("updates planet productions", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
		require.NotEqual(t, 9874, planet.Productions[0].Production)
		require.Nil(t, planet.Productions[0].Building)
		planet.Productions[0].Production = 9874
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs:        []models.BuildingActionCost{},
			Storages:     []models.BuildingActionResourceStorage{},
			// This should probably reflect the production of 9874 set with the planet
			// but as this does not play a role in the test it is left out
			Productions: []models.BuildingActionResourceProduction{},
		}
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, 9874)
	})

	t.Run("updates planet production for building", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
		require.NotEqual(t, 9874, planet.Productions[0].Production)
		require.NotNil(t, planet.Productions[0].Building)
		planet.Productions[0].Production = 9874
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     crystalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs:        []models.BuildingActionCost{},
			Storages:     []models.BuildingActionResourceStorage{},
			// This should probably reflect the value of 9874 set with the planet
			// but as this does not play a role in the test it is left out
			Productions: []models.BuildingActionResourceProduction{},
		}
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		assertPlanetResourceProduction(
			t, conn, planet.Id, metalResourceId, &crystalMineId, 9874,
		)
	})

	t.Run("does nothing when planet has no building action", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.Version++
		planet.BuildingAction = nil

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when planet has not the expected version", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.LessOrEqual(t, 1000.0, planet.Resources[0].Amount)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs: []models.BuildingActionCost{
				{
					Resource: crystalResourceId,
					Amount:   1000,
				},
			},
			Storages:    []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{},
		}

		// Only the resources are updated, the version stays the same:
		// this is not correct as the repository expects the version to
		// have been bumped.
		initialAmount := planet.Resources[0].Amount
		planet.Resources[0].Amount -= 1000

		err := repo.Create(t.Context(), planet)

		assert.ErrorIs(t, domainerrors.ErrOptimisticLocking, err, "Actual err: %v", err)
		assertBuildingActionDoesNotExist(t, conn, actionId)
		assertPlanetResourceAmount(t, conn, planet.Id, crystalResourceId, initialAmount)
	})

	t.Run("returns error when action for same planet already exists", func(t *testing.T) {
		_, planet := insertTestBuildingAction(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           uuid.New(),
			Building:     metalMineId,
			DesiredLevel: 5,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Minute),
		}

		err := repo.Create(t.Context(), planet)

		assert.Equal(t, domainerrors.ErrActionAlreadyInProgress, err, "Actual err: %v", err)
		assertBuildingActionDoesNotExist(t, conn, planet.BuildingAction.Id)
	})
}

func TestIT_BuildingActionRepository_Delete(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)

	t.Run("deletes an action", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn)
		planet.Version++

		planet.BuildingAction = nil
		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with costs", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionCost)
		planet.Version++

		planet.BuildingAction = nil
		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with storages", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionStorage)
		planet.Version++

		planet.BuildingAction = nil
		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with productions", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionProduction)
		planet.Version++

		planet.BuildingAction = nil
		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, action.Id)
	})

	t.Run("updates planet version and updated at", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		planet.BuildingAction = &action

		newTime := time.Date(2026, time.June, 20, 15, 2, 25, 0, time.UTC)
		require.NotEqual(t, planet.UpdatedAt, newTime)
		planet.UpdatedAt = newTime
		planet.Version++

		planet.BuildingAction = nil
		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)

		planetRepo := NewPlanetRepository(conn)
		actualPlanet, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet.Version, actualPlanet.Version)
		assert.Equal(t, newTime, actualPlanet.UpdatedAt)
	})

	t.Run("updates planet resource", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		planet.BuildingAction = &action

		planet.Resources[0].Amount += 1000
		planet.Version++

		planet.BuildingAction = nil
		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertPlanetResourceAmount(
			t,
			conn,
			planet.Id,
			planet.Resources[0].Resource,
			planet.Resources[0].Amount,
		)
	})

	t.Run("updates planet storages", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		planet.BuildingAction = &action

		require.NotEqual(t, planet.Storages[0].Storage, 5000)
		planet.Storages[0].Storage = 5000
		planet.Version++

		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, actionId)
		assertPlanetResourceStorage(t, conn, planet.Id, crystalResourceId, 5000)
	})

	t.Run("updates planet productions", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		planet.BuildingAction = &action

		require.NotEqual(t, 9874, planet.Productions[0].Production)
		require.Nil(t, planet.Productions[0].Building)
		planet.Productions[0].Production = 9874
		planet.Version++

		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, actionId)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, 9874)
	})

	t.Run("updates planet production for building", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		planet.BuildingAction = &action

		require.NotEqual(t, 9874, planet.Productions[0].Production)
		require.NotNil(t, planet.Productions[0].Building)
		planet.Productions[0].Production = 9874
		planet.Version++

		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, actionId)
		assertPlanetResourceProduction(
			t, conn, planet.Id, metalResourceId, &crystalMineId, 9874,
		)
	})
}

func TestIT_BuildingActionRepository_CreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	planetRepo := NewPlanetRepository(conn)

	type testCase struct {
		name   string
		action models.BuildingAction
	}

	testCases := []testCase{
		{
			name: "simple action",
			action: models.BuildingAction{
				Id:           uuid.New(),
				Building:     metalMineId,
				DesiredLevel: 27,
				CreatedAt:    time.Date(2024, time.December, 7, 20, 26, 47, 0, time.UTC),
				CompletedAt:  time.Date(2024, time.December, 7, 21, 26, 47, 0, time.UTC),
				Costs:        []models.BuildingActionCost{},
				Storages:     []models.BuildingActionResourceStorage{},
				Productions:  []models.BuildingActionResourceProduction{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			planet, _, _ := insertTestPlanetForPlayer(t, conn)
			planet.BuildingAction = &tc.action
			planet.Version++

			func() {
				err := repo.Create(t.Context(), planet)
				require.NoError(t, err, "Actual err: %v", err)
			}()

			func() {
				actual, err := planetRepo.Get(t.Context(), planet.Id)
				require.NoError(t, err, "Actual err: %v", err)

				require.NotNil(t, actual.BuildingAction)
				assert.Equal(t, tc.action, *actual.BuildingAction)
			}()

			planet.BuildingAction = nil
			planet.Version++

			func() {
				err := repo.Delete(t.Context(), planet)
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

func insertTestBuildingActionForPlanet(
	t *testing.T,
	conn db.Connection,
	planetId uuid.UUID,
	modifiers ...func(*testing.T, db.Connection, *models.BuildingAction),
) models.BuildingAction {
	t.Helper()

	action := models.BuildingAction{
		Id:           uuid.New(),
		Building:     metalMineId,
		DesiredLevel: 5,
		CreatedAt:    someTime,
		CompletedAt:  someTime.Add(1*time.Hour + 2*time.Minute),
		// This is intentional: the details (e.g. costs, productions, etc.) are returned as empty
		// slices by the adapter
		Costs:       []models.BuildingActionCost{},
		Storages:    []models.BuildingActionResourceStorage{},
		Productions: []models.BuildingActionResourceProduction{},
	}

	sqlQuery := `INSERT INTO building_action
		(id, planet, building, desired_level, created_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		action.Id,
		planetId,
		action.Building,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &action)
	}

	return action
}

func insertTestBuildingAction(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.BuildingAction),
) (models.BuildingAction, models.Planet) {
	t.Helper()

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	action := insertTestBuildingActionForPlanet(t, conn, planet.Id, modifiers...)
	planet.BuildingAction = &action
	return action, planet
}

func addBuildingActionCost(t *testing.T, conn db.Connection, a *models.BuildingAction) {
	t.Helper()

	insertBuildingActionCost(t, conn, metalResourceId, a)
}

func insertBuildingActionCost(t *testing.T, conn db.Connection, resourceId uuid.UUID, a *models.BuildingAction) {
	cost := models.BuildingActionCost{
		Resource: resourceId,
		Amount:   rand.Intn(4589),
	}

	sqlQuery := `INSERT INTO building_action_cost (action, resource, amount)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
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
		t.Context(),
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
		t.Context(),
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
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, 1, value)
}

func assertBuildingActionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionCostDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_cost WHERE action = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionStorageDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_storage WHERE action = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionProductionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_production WHERE action = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetResourceAmount(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, amount float64) {
	t.Helper()

	sqlQuery := `SELECT amount FROM planet_resource WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[float64](t.Context(), conn, sqlQuery, planet, resource)
	require.NoError(t, err, "Actual err: %v", err)
	require.InDelta(t, amount, value, 0.00001)
}

func assertPlanetResourceStorage(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, storage int) {
	t.Helper()

	sqlQuery := `SELECT storage FROM planet_resource_storage WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, resource)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, storage, value)
}

func assertPlanetResourceProduction(
	t *testing.T,
	conn db.Connection,
	planet uuid.UUID,
	resource uuid.UUID,
	building *uuid.UUID,
	production int,
) {
	t.Helper()

	sqlQuery := `SELECT production
		FROM planet_resource_production
		WHERE
			planet = $1
			AND resource = $2
			AND building IS NOT DISTINCT FROM $3`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, resource, building)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, production, value)
}
