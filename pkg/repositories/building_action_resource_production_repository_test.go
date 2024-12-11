package repositories

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/db/pgx"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	eassert "github.com/KnoblauchPilze/easy-assert/assert"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_BuildingActionResourceProductionRepository_Create(t *testing.T) {
	repo, conn, tx := newTestBuildingActionResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	action, _ := insertTestBuildingAction(t, conn)
	resource := insertTestResource(t, conn)

	production := persistence.BuildingActionResourceProduction{
		Action:     action.Id,
		Resource:   resource.Id,
		Production: 191,
	}

	actual, err := repo.Create(context.Background(), tx, production)
	assert.Nil(t, err)
	tx.Close(context.Background())

	assert.Equal(t, actual, production)
	assertBuildingActionResourceProductionExists(t, conn, action.Id, resource.Id)
	assertBuildingActionResourceProductionForResource(t, conn, action.Id, resource.Id, 191)
}

func TestIT_BuildingActionResourceProductionRepository_Create_WhenDuplicatedResource_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestBuildingActionResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	production, action, resource := insertTestBuildingActionResourceProduction(t, conn)

	newProduction := persistence.BuildingActionResourceProduction{
		Action:     action.Id,
		Resource:   resource.Id,
		Production: 34,
	}

	_, err := repo.Create(context.Background(), tx, newProduction)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertBuildingActionResourceProductionForResource(t, conn, action.Id, resource.Id, production.Production)
}

func TestIT_BuildingActionResourceProductionRepository_ListForAction(t *testing.T) {
	repo, conn, tx := newTestBuildingActionResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	barp1, action1, _ := insertTestBuildingActionResourceProduction(t, conn)
	barp2, _ := insertTestBuildingActionResourceProductionForAction(t, conn, action1.Id)
	barp3, action2, _ := insertTestBuildingActionResourceProduction(t, conn)

	actual, err := repo.ListForAction(context.Background(), tx, action1.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, barp1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, barp2))
	for _, buildingActionResourceProduction := range actual {
		assert.NotEqual(t, buildingActionResourceProduction.Action, action2.Id)
		assert.NotEqual(t, buildingActionResourceProduction.Resource, barp3.Resource)
	}
}

func newTestBuildingActionResourceProductionRepository(t *testing.T) (BuildingActionResourceProductionRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewBuildingActionResourceProductionRepository(), conn
}

func newTestBuildingActionResourceProductionRepositoryAndTransaction(t *testing.T) (BuildingActionResourceProductionRepository, db.Connection, db.Transaction) {
	repo, conn := newTestBuildingActionResourceProductionRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestBuildingActionResourceProduction(t *testing.T, conn db.Connection) (persistence.BuildingActionResourceProduction, persistence.BuildingAction, persistence.Resource) {
	action, _ := insertTestBuildingAction(t, conn)
	production, resource := insertTestBuildingActionResourceProductionForAction(t, conn, action.Id)
	return production, action, resource
}

func insertTestBuildingActionResourceProductionForAction(t *testing.T, conn db.Connection, action uuid.UUID) (persistence.BuildingActionResourceProduction, persistence.Resource) {
	resource := insertTestResource(t, conn)

	production := persistence.BuildingActionResourceProduction{
		Action:     action,
		Resource:   resource.Id,
		Production: 741,
	}

	sqlQuery := `INSERT INTO building_action_resource_production (action, resource, production) VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		production.Action,
		production.Resource,
		production.Production,
	)
	require.Nil(t, err)

	return production, resource
}

func assertBuildingActionResourceProductionExists(t *testing.T, conn db.Connection, action uuid.UUID, resource uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_production WHERE action = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action, resource)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertBuildingActionResourceProductionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_production WHERE action = $1`
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
