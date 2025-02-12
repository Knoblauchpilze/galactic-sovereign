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

func TestIT_BuildingActionResourceStorageRepository_Create(t *testing.T) {
	repo, conn, tx := newTestBuildingActionResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	action, _ := insertTestBuildingAction(t, conn)
	resource := insertTestResource(t, conn)

	storage := persistence.BuildingActionResourceStorage{
		Action:   action.Id,
		Resource: resource.Id,
		Storage:  879,
	}

	actual, err := repo.Create(context.Background(), tx, storage)
	assert.Nil(t, err)
	tx.Close(context.Background())

	assert.Equal(t, actual, storage)
	assertBuildingActionResourceStorageExists(t, conn, action.Id, resource.Id)
	assertBuildingActionResourceStorageForResource(t, conn, action.Id, resource.Id, 879)
}

func TestIT_BuildingActionResourceStorageRepository_Create_WhenDuplicatedResource_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestBuildingActionResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	storage, action, resource := insertTestBuildingActionResourceStorage(t, conn)

	newStorage := persistence.BuildingActionResourceStorage{
		Action:   action.Id,
		Resource: resource.Id,
		Storage:  storage.Storage + 10,
	}

	_, err := repo.Create(context.Background(), tx, newStorage)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertBuildingActionResourceStorageForResource(t, conn, action.Id, resource.Id, storage.Storage)
}

func TestIT_BuildingActionResourceStorageRepository_ListForAction(t *testing.T) {
	repo, conn, tx := newTestBuildingActionResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	bars1, action1, _ := insertTestBuildingActionResourceStorage(t, conn)
	bars2, _ := insertTestBuildingActionResourceStorageForAction(t, conn, action1.Id)
	bars3, action2, _ := insertTestBuildingActionResourceStorage(t, conn)

	actual, err := repo.ListForAction(context.Background(), tx, action1.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, bars1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, bars2))
	for _, buildingActionResourceStorage := range actual {
		assert.NotEqual(t, buildingActionResourceStorage.Action, action2.Id)
		assert.NotEqual(t, buildingActionResourceStorage.Resource, bars3.Resource)
	}
}

func newTestBuildingActionResourceStorageRepository(t *testing.T) (BuildingActionResourceStorageRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewBuildingActionResourceStorageRepository(), conn
}

func newTestBuildingActionResourceStorageRepositoryAndTransaction(t *testing.T) (BuildingActionResourceStorageRepository, db.Connection, db.Transaction) {
	repo, conn := newTestBuildingActionResourceStorageRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestBuildingActionResourceStorage(t *testing.T, conn db.Connection) (persistence.BuildingActionResourceStorage, persistence.BuildingAction, persistence.Resource) {
	action, _ := insertTestBuildingAction(t, conn)
	storage, resource := insertTestBuildingActionResourceStorageForAction(t, conn, action.Id)
	return storage, action, resource
}

func insertTestBuildingActionResourceStorageForAction(t *testing.T, conn db.Connection, action uuid.UUID) (persistence.BuildingActionResourceStorage, persistence.Resource) {
	resource := insertTestResource(t, conn)

	storage := persistence.BuildingActionResourceStorage{
		Action:   action,
		Resource: resource.Id,
		Storage:  546,
	}

	sqlQuery := `INSERT INTO building_action_resource_storage (action, resource, storage) VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		storage.Action,
		storage.Resource,
		storage.Storage,
	)
	require.Nil(t, err)

	return storage, resource
}

func assertBuildingActionResourceStorageExists(t *testing.T, conn db.Connection, action uuid.UUID, resource uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_storage WHERE action = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action, resource)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertBuildingActionResourceStorageForResource(t *testing.T, conn db.Connection, action uuid.UUID, resource uuid.UUID, storage int) {
	sqlQuery := `SELECT storage FROM building_action_resource_storage WHERE action = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action, resource)
	require.Nil(t, err)
	require.Equal(t, storage, value)
}
