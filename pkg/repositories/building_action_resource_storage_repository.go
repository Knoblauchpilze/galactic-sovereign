package repositories

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionResourceStorageRepository interface {
	Create(ctx context.Context, tx db.Transaction, storage persistence.BuildingActionResourceStorage) (persistence.BuildingActionResourceStorage, error)
	ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionResourceStorage, error)
}

type buildingActionResourceStorageRepositoryImpl struct{}

func NewBuildingActionResourceStorageRepository() BuildingActionResourceStorageRepository {
	return &buildingActionResourceStorageRepositoryImpl{}
}

const createBuildingActionResourceStorageSqlTemplate = `
INSERT INTO
	building_action_resource_storage (action, resource, storage)
	VALUES($1, $2, $3)`

func (r *buildingActionResourceStorageRepositoryImpl) Create(ctx context.Context, tx db.Transaction, storage persistence.BuildingActionResourceStorage) (persistence.BuildingActionResourceStorage, error) {
	_, err := tx.Exec(ctx, createBuildingActionResourceStorageSqlTemplate, storage.Action, storage.Resource, storage.Storage)
	return storage, err
}

const listBuildingActionResourceStorageForActionSqlTemplate = `
SELECT
	action,
	resource,
	storage
FROM
	building_action_resource_storage
WHERE
	action = $1`

func (r *buildingActionResourceStorageRepositoryImpl) ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionResourceStorage, error) {
	return db.QueryAllTx[persistence.BuildingActionResourceStorage](ctx, tx, listBuildingActionResourceStorageForActionSqlTemplate, action)
}
