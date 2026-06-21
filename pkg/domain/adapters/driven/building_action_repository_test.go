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

	t.Run("creates an action", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
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
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := repo.Get(t.Context(), actionId)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, *planet.BuildingAction, actual)
	})

	t.Run("creates an action with costs", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
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
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := repo.Get(t.Context(), actionId)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, *planet.BuildingAction, actual)
	})

	t.Run("creates an action with storages", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
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
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := repo.Get(t.Context(), actionId)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, *planet.BuildingAction, actual)
	})

	t.Run("creates an action with productions", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
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
		planet.Version++

		err := repo.Create(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertBuildingActionExists(t, conn, actionId)

		actual, err := repo.Get(t.Context(), actionId)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, *planet.BuildingAction, actual)
	})

	t.Run("updates planet version and updated at", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Planet:       planet.Id,
			Building:     metalMineId,
			CurrentLevel: 2,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Version:      9,
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

		planetRepo := NewPlanetRepository(conn)
		actualPlanet, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet.Version, actualPlanet.Version)
		assert.Equal(t, newTime, actualPlanet.UpdatedAt)
	})

	t.Run("updates planet resources", func(t *testing.T) {
		actionId := uuid.New()
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		// This is to make sure that the cost is not bigger than the amount of
		// resources on the planet
		require.LessOrEqual(t, 1000.0, planet.Resources[0].Amount)
		planet.BuildingAction = &models.BuildingAction{
			Id:           actionId,
			Planet:       planet.Id,
			Building:     metalMineId,
			CurrentLevel: 2,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Version:      9,
			Costs: []models.BuildingActionCost{
				{
					Resource: crystalResourceId,
					Amount:   1000,
				},
			},
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

	t.Run("returns error when action for same planet already exists", func(t *testing.T) {
		_, planet := insertTestBuildingAction(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			CurrentLevel: 4,
			DesiredLevel: 5,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Minute),
			Version:      14,
		}

		err := repo.Create(t.Context(), planet)

		assert.Equal(t, domainerrors.ErrActionAlreadyInProgress, err, "Actual err: %v", err)
		assertBuildingActionDoesNotExist(t, conn, planet.BuildingAction.Id)
	})
}

func TestIT_BuildingActionRepository_Get(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)

	t.Run("gets an action", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn)

		actual, err := repo.Get(t.Context(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("gets an action with costs", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionCost)

		actual, err := repo.Get(t.Context(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("gets an action with storages", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		actual, err := repo.Get(t.Context(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("gets an action with productions", func(t *testing.T) {
		action, _ := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		actual, err := repo.Get(t.Context(), action.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, action, actual)
	})

	t.Run("returns error when action does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func TestIT_BuildingActionRepository_Delete(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)

	t.Run("deletes an action", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn)
		planet.Version++

		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with costs", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionCost)
		planet.Version++

		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with storages", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionStorage)
		planet.Version++

		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, action.Id)
	})

	t.Run("deletes action with productions", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionProduction)
		planet.Version++

		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, action.Id)
	})

	t.Run("updates planet version and updated at", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		modifier := func(t *testing.T, conn db.Connection, a *models.BuildingAction) {
			insertBuildingActionCost(t, conn, planet.Resources[0].Resource, a)
		}
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, modifier)

		newTime := time.Date(2026, time.June, 20, 15, 2, 25, 0, time.UTC)
		require.NotEqual(t, planet.UpdatedAt, newTime)
		planet.UpdatedAt = newTime
		planet.Version++

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
		modifier := func(t *testing.T, conn db.Connection, a *models.BuildingAction) {
			insertBuildingActionCost(t, conn, planet.Resources[0].Resource, a)
		}
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, modifier)

		planet.Resources[0].Amount += 1000
		planet.Version++

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

	t.Run("succeeds when the action does not exist", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		planet.BuildingAction = &models.BuildingAction{
			Id:     uuid.New(),
			Planet: planet.Id,
		}

		planet.Version++

		err := repo.Delete(t.Context(), planet)
		require.NoError(t, err, "Actual err: %v", err)
	})
}

func TestIT_BuildingActionRepository_CreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)

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
				CurrentLevel: 26,
				DesiredLevel: 27,
				CreatedAt:    time.Date(2024, time.December, 7, 20, 26, 47, 0, time.UTC),
				CompletedAt:  time.Date(2024, time.December, 7, 21, 26, 47, 0, time.UTC),
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
			planet.BuildingAction = &tc.action
			planet.Version++

			func() {
				err := repo.Create(t.Context(), planet)
				require.NoError(t, err, "Actual err: %v", err)
			}()

			func() {
				actionFromDb, err := repo.Get(t.Context(), tc.action.Id)
				require.NoError(t, err, "Actual err: %v", err)

				assert.Equal(t, tc.action, actionFromDb)
			}()

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
		Planet:       planetId,
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
		t.Context(),
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
