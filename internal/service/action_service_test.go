package service

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/game"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_ActionService_ProcessActionUntil_WheNoAction_ExpectSuccess(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)

	veryFarInThePast := time.Now().Add(-800 * time.Hour)

	err := service.ProcessActionsUntil(context.Background(), planet.Id, veryFarInThePast)

	assert.Nil(t, err)
}

func TestIT_ActionService_ProcessActionUntil_WhenActionIsNotCompleted_ExpectNotProcessed(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	createdAt := time.Now().Add(-2 * time.Hour)
	completedAt := createdAt.Add(8 * time.Hour)
	action, _ := insertTestBuildingActionForPlanetWithTimes(
		t,
		conn,
		planet.Id,
		createdAt,
		completedAt,
	)
	beforeActionCompletes := completedAt.Add(-1 * time.Hour)

	err := service.ProcessActionsUntil(context.Background(), planet.Id, beforeActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionExists(t, conn, action.Id)
}

func TestIT_ActionService_ProcessActionUntil_WhenActionIsCompleted_ExpectRemoval(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	createdAt := time.Now().Add(-2 * time.Hour)
	completedAt := createdAt.Add(-1 * time.Hour)
	action := insertTestBuildingActionForPlanetAndBuildingWithTimes(
		t,
		conn,
		planet.Id,
		planetBuilding.Building,
		createdAt,
		completedAt,
	)
	afterActionCompletes := completedAt.Add(1 * time.Second)

	err := service.ProcessActionsUntil(context.Background(), planet.Id, afterActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
}

func TestIT_ActionService_ProcessActionUntil_ExpectResourcesToBeUpdatedToCompletionTime(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	planetResourceProd, res := insertTestPlanetResourceProductionForBuilding(t, conn, planet.Id, planetBuilding.Building)
	insertTestPlanetResourceStorageForResource(t, conn, planet.Id, res.Id, 50_000)
	createdAt := time.Now().Add(-3 * time.Hour)
	completedAt := createdAt.Add(2 * time.Hour)
	action := insertTestBuildingActionForPlanetAndBuildingWithTimes(
		t,
		conn,
		planet.Id,
		planetBuilding.Building,
		createdAt,
		completedAt,
	)
	updatedAt := createdAt.Add(1 * time.Hour)
	planetResource := insertTestPlanetResourceForResource(t, conn, planet.Id, res.Id, updatedAt)

	afterActionCompletes := completedAt.Add(1 * time.Second)
	err := service.ProcessActionsUntil(context.Background(), planet.Id, afterActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	elapsed := completedAt.Sub(updatedAt)
	expectedAmount := math.Floor(float64(planetResource.Amount) + elapsed.Hours()*float64(planetResourceProd.Production))
	assertPlanetResourceAmount(
		t,
		conn,
		planetResourceProd.Planet,
		planetResourceProd.Resource,
		expectedAmount,
	)
}

func TestIT_ActionService_ProcessActionUntil_WhenResourceIsNotProduced_ExpectResourceToStayTheSame(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	res := insertTestResource(t, conn)
	insertTestPlanetResourceStorageForResource(t, conn, planet.Id, res.Id, 50_000)
	createdAt := time.Now().Add(-3 * time.Hour)
	completedAt := createdAt.Add(2 * time.Hour)
	action := insertTestBuildingActionForPlanetAndBuildingWithTimes(
		t,
		conn,
		planet.Id,
		planetBuilding.Building,
		createdAt,
		completedAt,
	)
	updatedAt := createdAt.Add(1 * time.Hour)
	planetResource := insertTestPlanetResourceForResource(t, conn, planet.Id, res.Id, updatedAt)

	afterActionCompletes := completedAt.Add(1 * time.Second)
	err := service.ProcessActionsUntil(context.Background(), planet.Id, afterActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	assertPlanetResourceAmount(
		t,
		conn,
		planet.Id,
		res.Id,
		planetResource.Amount,
	)
}

func TestIT_ActionService_ProcessActionUntil_WhenStorageIsAlreadyFull_ExpectResourceToStayTheSame(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	planetResourceProd, res := insertTestPlanetResourceProductionForBuilding(t, conn, planet.Id, planetBuilding.Building)
	createdAt := time.Now().Add(-3 * time.Hour)
	completedAt := createdAt.Add(2 * time.Hour)
	action := insertTestBuildingActionForPlanetAndBuildingWithTimes(
		t,
		conn,
		planet.Id,
		planetBuilding.Building,
		createdAt,
		completedAt,
	)
	updatedAt := createdAt.Add(1 * time.Hour)
	planetResource := insertTestPlanetResourceForResource(t, conn, planet.Id, res.Id, updatedAt)
	lowerStorageThanWhatAlreadyExists := int(planetResource.Amount - 100)
	insertTestPlanetResourceStorageForResource(
		t,
		conn,
		planet.Id,
		res.Id,
		lowerStorageThanWhatAlreadyExists,
	)

	afterActionCompletes := completedAt.Add(1 * time.Second)
	err := service.ProcessActionsUntil(context.Background(), planet.Id, afterActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	assertPlanetResourceAmount(
		t,
		conn,
		planetResourceProd.Planet,
		planetResourceProd.Resource,
		planetResource.Amount,
	)
}

func TestIT_ActionService_ProcessActionUntil_WhenStorageCanNotHoldAllProduction_ExpectResourceToFillStorage(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	planetResourceProd, res := insertTestPlanetResourceProductionForBuilding(t, conn, planet.Id, planetBuilding.Building)
	createdAt := time.Now().Add(-3 * time.Hour)
	completedAt := createdAt.Add(2 * time.Hour)
	action := insertTestBuildingActionForPlanetAndBuildingWithTimes(
		t,
		conn,
		planet.Id,
		planetBuilding.Building,
		createdAt,
		completedAt,
	)
	updatedAt := createdAt.Add(1 * time.Hour)
	planetResource := insertTestPlanetResourceForResource(t, conn, planet.Id, res.Id, updatedAt)

	elapsed := completedAt.Sub(updatedAt)
	expectedProduction := elapsed.Hours() * float64(planetResourceProd.Production)
	fullAmount := float64(planetResource.Amount) + expectedProduction
	lowerStorageThanNeeded := fullAmount - 500
	insertTestPlanetResourceStorageForResource(
		t,
		conn,
		planet.Id,
		res.Id,
		int(lowerStorageThanNeeded),
	)

	afterActionCompletes := completedAt.Add(1 * time.Second)
	err := service.ProcessActionsUntil(context.Background(), planet.Id, afterActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	expectedAmount := lowerStorageThanNeeded
	assertPlanetResourceAmount(
		t,
		conn,
		planetResourceProd.Planet,
		planetResourceProd.Resource,
		expectedAmount,
	)
}

func TestIT_ActionService_ProcessActionUntil_ExpectBuildingToBeUpdated(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	createdAt := time.Now().Add(-3 * time.Hour)
	completedAt := createdAt.Add(2 * time.Hour)
	action := insertTestBuildingActionForPlanetAndBuildingWithTimes(
		t,
		conn,
		planet.Id,
		planetBuilding.Building,
		createdAt,
		completedAt,
	)

	afterActionCompletes := completedAt.Add(1 * time.Second)
	err := service.ProcessActionsUntil(context.Background(), planet.Id, afterActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	expectedLevel := action.DesiredLevel
	assertPlanetBuildingLevel(t, conn, planet.Id, planetBuilding.Building, expectedLevel)
}

func TestIT_ActionService_ProcessActionUntil_ExpectResourceProductionToBeUpdated(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	prp1, res1 := insertTestPlanetResourceProductionForBuilding(t, conn, planet.Id, planetBuilding.Building)
	prp2, res2 := insertTestPlanetResourceProduction(t, conn, planet.Id)
	createdAt := time.Now().Add(-3 * time.Hour)
	completedAt := createdAt.Add(2 * time.Hour)
	action := insertTestBuildingActionForPlanetAndBuildingWithTimes(
		t,
		conn,
		planet.Id,
		planetBuilding.Building,
		createdAt,
		completedAt,
	)
	actionProduction := insertTestBuildingActionResourceProductionForActionAndResource(
		t,
		conn,
		action.Id,
		res1.Id,
	)

	afterActionCompletes := completedAt.Add(1 * time.Second)
	err := service.ProcessActionsUntil(context.Background(), planet.Id, afterActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	assertPlanetResourceProductionForBuilding(
		t,
		conn,
		planet.Id,
		res1.Id,
		*prp1.Building,
		actionProduction.Production,
	)
	assertPlanetResourceProductionWithoutBuilding(
		t,
		conn,
		planet.Id,
		res2.Id,
		prp2.Production,
	)
}

func TestIT_ActionService_ProcessActionUntil_WhenResourceIsNotProduced_ExpectResourceProductionCreatedToNewValue(t *testing.T) {
	service, conn := newTestActionService(t)
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	res := insertTestResource(t, conn)
	createdAt := time.Now().Add(-3 * time.Hour)
	completedAt := createdAt.Add(2 * time.Hour)
	action := insertTestBuildingActionForPlanetAndBuildingWithTimes(
		t,
		conn,
		planet.Id,
		planetBuilding.Building,
		createdAt,
		completedAt,
	)
	actionProduction := insertTestBuildingActionResourceProductionForActionAndResource(
		t,
		conn,
		action.Id,
		res.Id,
	)

	afterActionCompletes := completedAt.Add(1 * time.Second)
	err := service.ProcessActionsUntil(context.Background(), planet.Id, afterActionCompletes)

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
	assertPlanetResourceProductionForBuilding(
		t,
		conn,
		planet.Id,
		res.Id,
		planetBuilding.Building,
		actionProduction.Production,
	)
}

func newTestActionService(t *testing.T) (game.ActionService, db.Connection) {
	conn := newTestConnection(t)
	repos := repositories.Repositories{
		BuildingAction:                   repositories.NewBuildingActionRepository(),
		BuildingActionResourceProduction: repositories.NewBuildingActionResourceProductionRepository(),
		BuildingActionResourceStorage:    repositories.NewBuildingActionResourceStorageRepository(),
		PlanetBuilding:                   repositories.NewPlanetBuildingRepository(),
		PlanetResource:                   repositories.NewPlanetResourceRepository(),
		PlanetResourceProduction:         repositories.NewPlanetResourceProductionRepository(),
		PlanetResourceStorage:            repositories.NewPlanetResourceStorageRepository(),
	}

	service := NewActionService(conn, repos)

	return service, conn
}

func assertPlanetBuildingLevel(t *testing.T, conn db.Connection, planet uuid.UUID, building uuid.UUID, level int) {
	sqlQuery := `SELECT level FROM planet_building WHERE planet = $1 AND building = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, building)
	require.Nil(t, err)
	require.Equal(t, level, value)
}

func assertPlanetResourceProductionForBuilding(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, building uuid.UUID, prod int) {
	sqlQuery := `SELECT production FROM planet_resource_production WHERE planet = $1 AND resource = $2 AND building = $3`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, resource, building)
	require.Nil(t, err)
	require.Equal(t, prod, value)
}

func assertPlanetResourceProductionWithoutBuilding(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, prod int) {
	sqlQuery := `SELECT production FROM planet_resource_production WHERE planet = $1 AND resource = $2 AND building IS NULL`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, prod, value)
}
