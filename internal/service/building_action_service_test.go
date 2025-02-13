package service

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/game"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_BuildingActionService_Create(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Building,
	}
	actionResponse, err := service.Create(context.Background(), actionRequest)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionExists(t, conn, actionResponse.Id)
	expected := communication.BuildingActionDtoResponse{
		Planet:       actionRequest.Planet,
		Building:     actionRequest.Building,
		CurrentLevel: building.Level,
		DesiredLevel: building.Level + 1,
	}
	assert.True(t, eassert.EqualsIgnoringFields(actionResponse, expected, "Id", "CreatedAt", "CompletedAt"))
}

func TestIT_BuildingActionService_Create_WhenNotEnoughResources_ExpectFailure(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	insertTestBuildingCost(t, conn, building.Building)

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Building,
	}
	_, err := service.Create(context.Background(), actionRequest)

	assert.True(t, errors.IsErrorWithCode(err, game.NotEnoughResources))
}

func TestIT_BuildingActionService_Create_WhenBuildingDoesNotExist_ExpectFailure(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building := insertTestBuilding(t, conn)
	insertTestBuildingCost(t, conn, building.Id)

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Id,
	}
	_, err := service.Create(context.Background(), actionRequest)

	assert.True(t, errors.IsErrorWithCode(err, game.NoSuchBuilding))
}

func TestIT_BuildingActionService_Create_WithCost_ExpectCostToBeRegistered(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	cost, _ := insertTestBuildingCost(t, conn, building.Building)
	insertTestPlanetResourceForResource(t, conn, planet.Id, cost.Resource, time.Now())

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Building,
	}
	actionResponse, err := service.Create(context.Background(), actionRequest)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionExists(t, conn, actionResponse.Id)
	assertBuildingActionCostForResource(t, conn, actionResponse.Id, cost.Resource, 268)
}

func TestIT_BuildingActionService_Create_WithCost_ExpectCostToBeTakenFromResources(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	cost, _ := insertTestBuildingCost(t, conn, building.Building)
	planetResource := insertTestPlanetResourceForResource(
		t, conn, planet.Id, cost.Resource, time.Now(),
	)

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Building,
	}
	actionResponse, err := service.Create(context.Background(), actionRequest)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionExists(t, conn, actionResponse.Id)
	expectedCost := 268.0
	assertPlanetResourceAmount(t, conn, planet.Id, cost.Resource, planetResource.Amount-expectedCost)
}

func TestIT_BuildingActionService_Create_WithResourceProduction_ExpectProductionToBeRegistered(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	production, _ := insertTestBuildingResourceProduction(t, conn, building.Building)

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Building,
	}
	actionResponse, err := service.Create(context.Background(), actionRequest)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionExists(t, conn, actionResponse.Id)
	assertBuildingActionResourceProductionForResource(
		t,
		conn,
		actionResponse.Id,
		production.Resource,
		12_176_686,
	)
}

func TestIT_BuildingActionService_Create_WithResourceStorage_ExpectStorageToBeRegistered(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	storage, _ := insertTestBuildingResourceStorage(t, conn, building.Building)

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Building,
	}
	actionResponse, err := service.Create(context.Background(), actionRequest)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionExists(t, conn, actionResponse.Id)
	assertBuildingActionResourceStorageForResource(
		t,
		conn,
		actionResponse.Id,
		storage.Resource,
		359_260_928,
	)
}

func TestIT_BuildingActionService_Delete(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	createdAt := time.Now()
	completedAt := createdAt.Add(2 * time.Hour)
	action, _ := insertTestBuildingActionForPlanetWithTimes(
		t,
		conn,
		planet.Id,
		createdAt,
		completedAt,
	)

	err := service.Delete(context.Background(), action.Id)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
}

func TestIT_BuildingActionService_Delete_ExpectCostsToBeDeleted(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	createdAt := time.Now()
	completedAt := createdAt.Add(2 * time.Hour)
	action, _ := insertTestBuildingActionForPlanetWithTimes(t, conn, planet.Id, createdAt, completedAt)
	_, res := insertTestBuildingActionCostForAction(t, conn, action.Id)
	insertTestPlanetResourceForResource(t, conn, planet.Id, res.Id, time.Now())

	err := service.Delete(context.Background(), action.Id)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	assertBuildingActionCostDoesNotExist(t, conn, action.Id)
}

func TestIT_BuildingActionService_Delete_ExpectResourcesToBeRestored(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	createdAt := time.Now()
	completedAt := createdAt.Add(2 * time.Hour)
	action, _ := insertTestBuildingActionForPlanetWithTimes(t, conn, planet.Id, createdAt, completedAt)
	cost, _ := insertTestBuildingActionCostForAction(t, conn, action.Id)
	planetResource := insertTestPlanetResourceForResource(
		t, conn, planet.Id, cost.Resource, time.Now(),
	)

	err := service.Delete(context.Background(), action.Id)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	expectedAmount := planetResource.Amount + float64(cost.Amount)
	assertPlanetResourceAmount(t, conn, planet.Id, planetResource.Resource, expectedAmount)
}

func TestIT_BuildingActionService_Delete_ExpectResourceProductionToBeDeleted(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	createdAt := time.Now()
	completedAt := createdAt.Add(2 * time.Hour)
	action, _ := insertTestBuildingActionForPlanetWithTimes(t, conn, planet.Id, createdAt, completedAt)
	insertTestBuildingActionResourceProductionForAction(t, conn, action.Id)

	err := service.Delete(context.Background(), action.Id)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	assertBuildingActionResourceProductionDoesNotExist(t, conn, action.Id)
}

func TestIT_BuildingActionService_Delete_ExpectResourceStorageToBeDeleted(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	createdAt := time.Now()
	completedAt := createdAt.Add(2 * time.Hour)
	action, _ := insertTestBuildingActionForPlanetWithTimes(t, conn, planet.Id, createdAt, completedAt)
	insertTestBuildingActionResourceStorageForAction(t, conn, action.Id)

	err := service.Delete(context.Background(), action.Id)

	assert.Nil(t, err, "Actual err: %v", err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	assertBuildingActionResourceStorageDoesNotExist(t, conn, action.Id)
}

func TestIT_BuildingActionService_Delete_WhenActionAlreadyCompleted_ExpectFailure(t *testing.T) {
	service, conn := newTestBuildingActionService(t)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	createdAt := time.Now().Add(-23 * time.Hour)
	completedAt := createdAt.Add(30 * time.Minute)
	action, _ := insertTestBuildingActionForPlanetWithTimes(t, conn, planet.Id, createdAt, completedAt)

	err := service.Delete(context.Background(), action.Id)

	assert.True(t, errors.IsErrorWithCode(err, ActionAlreadyCompleted), "Actual err: %v", err)
}

func TestIT_BuildingActionService_CreationDeletionWorkflow(t *testing.T) {
	var generatedCreatedAt time.Time
	var returnedCompletionTime time.Time
	completionTimeFunc := func(action persistence.BuildingAction, resources []persistence.Resource, costs []persistence.BuildingActionCost) (persistence.BuildingAction, error) {
		generatedCreatedAt = action.CreatedAt
		returnedCompletionTime = time.Now().Add(1 * time.Hour)
		action.CompletedAt = returnedCompletionTime
		return action, nil
	}

	service, conn := newTestBuildingActionServiceWithCompletionTime(t, completionTimeFunc)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	_, resource := insertTestBuildingCost(t, conn, building.Building)
	insertTestPlanetResourceForResource(t, conn, planet.Id, resource.Id, time.Now())

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Building,
	}

	var err error
	var actionResponse communication.BuildingActionDtoResponse
	func() {
		actionResponse, err = service.Create(context.Background(), actionRequest)
		require.Nil(t, err)
	}()

	assertBuildingActionExists(t, conn, actionResponse.Id)
	expected := communication.BuildingActionDtoResponse{
		Planet:       actionRequest.Planet,
		Building:     actionRequest.Building,
		CurrentLevel: 4,
		DesiredLevel: 5,
		CreatedAt:    generatedCreatedAt,
		CompletedAt:  returnedCompletionTime,
	}
	assert.True(t, eassert.EqualsIgnoringFields(actionResponse, expected, "Id"))

	func() {
		err = service.Delete(context.Background(), actionResponse.Id)
		require.Nil(t, err)
	}()

	assertBuildingActionDoesNotExist(t, conn, actionResponse.Id)
}

func newTestBuildingActionService(t *testing.T) (BuildingActionService, db.Connection) {
	conn := newTestConnection(t)
	repos := repositories.Repositories{
		Resource:                         repositories.NewResourceRepository(),
		PlanetResource:                   repositories.NewPlanetResourceRepository(),
		PlanetBuilding:                   repositories.NewPlanetBuildingRepository(),
		BuildingCost:                     repositories.NewBuildingCostRepository(),
		BuildingResourceProduction:       repositories.NewBuildingResourceProductionRepository(),
		BuildingResourceStorage:          repositories.NewBuildingResourceStorageRepository(),
		BuildingAction:                   repositories.NewBuildingActionRepository(),
		BuildingActionCost:               repositories.NewBuildingActionCostRepository(),
		BuildingActionResourceProduction: repositories.NewBuildingActionResourceProductionRepository(),
		BuildingActionResourceStorage:    repositories.NewBuildingActionResourceStorageRepository(),
	}

	service := NewBuildingActionService(conn, repos)

	return service, conn
}

func newTestBuildingActionServiceWithCompletionTime(t *testing.T, consolidator buildingActionCompletionTimeConsolidator) (BuildingActionService, db.Connection) {
	conn := newTestConnection(t)
	repos := repositories.Repositories{
		Resource:                         repositories.NewResourceRepository(),
		PlanetResource:                   repositories.NewPlanetResourceRepository(),
		PlanetBuilding:                   repositories.NewPlanetBuildingRepository(),
		BuildingCost:                     repositories.NewBuildingCostRepository(),
		BuildingResourceProduction:       repositories.NewBuildingResourceProductionRepository(),
		BuildingResourceStorage:          repositories.NewBuildingResourceStorageRepository(),
		BuildingAction:                   repositories.NewBuildingActionRepository(),
		BuildingActionCost:               repositories.NewBuildingActionCostRepository(),
		BuildingActionResourceProduction: repositories.NewBuildingActionResourceProductionRepository(),
		BuildingActionResourceStorage:    repositories.NewBuildingActionResourceStorageRepository(),
	}

	service := newBuildingActionServiceWithCompletionTime(conn, repos, consolidator)

	return service, conn
}

func assertBuildingActionExists(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertBuildingActionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertBuildingActionCostForResource(t *testing.T, conn db.Connection, action uuid.UUID, resource uuid.UUID, cost int) {
	sqlQuery := `SELECT amount FROM building_action_cost WHERE action = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action, resource)
	require.Nil(t, err)
	require.Equal(t, cost, value)
}

func assertBuildingActionCostDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action_cost WHERE action = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertBuildingActionResourceProductionForResource(t *testing.T, conn db.Connection, action uuid.UUID, resource uuid.UUID, production int) {
	sqlQuery := `SELECT production FROM building_action_resource_production WHERE action = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action, resource)
	require.Nil(t, err)
	require.Equal(t, production, value)
}

func assertBuildingActionResourceProductionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_production WHERE action = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertBuildingActionResourceStorageForResource(t *testing.T, conn db.Connection, action uuid.UUID, resource uuid.UUID, storage int) {
	sqlQuery := `SELECT storage FROM building_action_resource_storage WHERE action = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action, resource)
	require.Nil(t, err)
	require.Equal(t, storage, value)
}

func assertBuildingActionResourceStorageDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_storage WHERE action = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertPlanetResourceAmount(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, amount float64) {
	sqlQuery := `SELECT amount FROM planet_resource WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[float64](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, amount, value)
}
