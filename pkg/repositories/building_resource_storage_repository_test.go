package repositories

import (
	"context"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_BuildingResourceStorageRepository_List(t *testing.T) {
	repo, conn, tx := newTestBuildingResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	brs1, b, _ := insertTestBuildingResourceStorageForBuilding(t, conn)
	brs2, _ := insertTestBuildingResourceStorage(t, conn, b.Id)

	actual, err := repo.ListForBuilding(context.Background(), tx, b.Id)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, brs1)
	assert.Contains(t, actual, brs2)
}

func newTestBuildingResourceStorageRepositoryAndTransaction(t *testing.T) (BuildingResourceStorageRepository, db.Connection, db.Transaction) {
	conn := newTestConnection(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return NewBuildingResourceStorageRepository(), conn, tx
}

func insertTestBuildingResourceStorage(t *testing.T, conn db.Connection, building uuid.UUID) (persistence.BuildingResourceStorage, persistence.Resource) {
	resource := insertTestResource(t, conn)

	storage := persistence.BuildingResourceStorage{
		Building: building,
		Resource: resource.Id,
		Base:     89,
		Scale:    8.12,
		Progress: 1.15,
	}

	sqlQuery := `INSERT INTO building_resource_storage (building, resource, base, scale, progress) VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		storage.Building,
		storage.Resource,
		storage.Base,
		storage.Scale,
		storage.Progress,
	)
	require.Nil(t, err)

	return storage, resource
}

func insertTestBuildingResourceStorageForBuilding(t *testing.T, conn db.Connection) (persistence.BuildingResourceStorage, persistence.Building, persistence.Resource) {
	building := insertTestBuilding(t, conn)
	storage, resource := insertTestBuildingResourceStorage(t, conn, building.Id)
	return storage, building, resource
}
