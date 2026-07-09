package drivenadapters

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	metalMineId = uuid.MustParse("d176e82d-f2ca-4611-996b-c4804096caef")
)

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
