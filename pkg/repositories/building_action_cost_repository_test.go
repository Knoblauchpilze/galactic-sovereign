package repositories

import (
	"context"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_BuildingActionCostRepository_Create(t *testing.T) {
	repo, conn, tx := newTestBuildingActionCostRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	action, _ := insertTestBuildingAction(t, conn)
	resource := insertTestResource(t, conn)

	cost := persistence.BuildingActionCost{
		Action:   action.Id,
		Resource: resource.Id,
		Amount:   26,
	}

	actual, err := repo.Create(context.Background(), tx, cost)
	assert.Nil(t, err)
	tx.Close(context.Background())

	assert.Equal(t, actual, cost)
	assertBuildingActionCostExists(t, conn, action.Id, resource.Id)
	assertBuildingActionCostForResource(t, conn, action.Id, resource.Id, 26)
}

func TestIT_BuildingActionCostRepository_Create_WhenDuplicatedResource_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestBuildingActionCostRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	cost, action, resource := insertTestBuildingActionCost(t, conn)

	newCost := persistence.BuildingActionCost{
		Action:   action.Id,
		Resource: resource.Id,
		Amount:   58,
	}

	_, err := repo.Create(context.Background(), tx, newCost)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertBuildingActionCostForResource(t, conn, action.Id, resource.Id, cost.Amount)
}

func TestIT_BuildingActionCostRepository_ListForAction(t *testing.T) {
	repo, conn, tx := newTestBuildingActionCostRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	bac1, action1, _ := insertTestBuildingActionCost(t, conn)
	bac2, _ := insertTestBuildingActionCostForAction(t, conn, action1.Id)
	bac3, action2, _ := insertTestBuildingActionCost(t, conn)

	actual, err := repo.ListForAction(context.Background(), tx, action1.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, bac1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, bac2))
	for _, buildingActionCost := range actual {
		assert.NotEqual(t, buildingActionCost.Action, action2.Id)
		assert.NotEqual(t, buildingActionCost.Resource, bac3.Resource)
	}
}

func newTestBuildingActionCostRepository(t *testing.T) (BuildingActionCostRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewBuildingActionCostRepository(), conn
}

func newTestBuildingActionCostRepositoryAndTransaction(t *testing.T) (BuildingActionCostRepository, db.Connection, db.Transaction) {
	repo, conn := newTestBuildingActionCostRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestBuildingActionCost(t *testing.T, conn db.Connection) (persistence.BuildingActionCost, persistence.BuildingAction, persistence.Resource) {
	action, _ := insertTestBuildingAction(t, conn)
	cost, resource := insertTestBuildingActionCostForAction(t, conn, action.Id)
	return cost, action, resource
}

func insertTestBuildingActionCostForAction(t *testing.T, conn db.Connection, action uuid.UUID) (persistence.BuildingActionCost, persistence.Resource) {
	resource := insertTestResource(t, conn)

	cost := persistence.BuildingActionCost{
		Action:   action,
		Resource: resource.Id,
		Amount:   56,
	}

	sqlQuery := `INSERT INTO building_action_cost (action, resource, amount) VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		cost.Action,
		cost.Resource,
		cost.Amount,
	)
	require.Nil(t, err)

	return cost, resource
}

func assertBuildingActionCostExists(t *testing.T, conn db.Connection, action uuid.UUID, resource uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action_cost WHERE action = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action, resource)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertBuildingActionCostDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action_cost WHERE action = $1`
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
