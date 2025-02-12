package repositories

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingResourceStorageRepository interface {
	ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingResourceStorage, error)
}

type buildingResourceStorageRepositoryImpl struct{}

func NewBuildingResourceStorageRepository() BuildingResourceStorageRepository {
	return &buildingResourceStorageRepositoryImpl{}
}

const listBuildingResourceStorageForBuildingSqlTemplate = `
SELECT
	building,
	resource,
	base,
	scale,
	progress
FROM
	building_resource_storage
WHERE
	building = $1`

func (r *buildingResourceStorageRepositoryImpl) ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingResourceStorage, error) {
	return db.QueryAllTx[persistence.BuildingResourceStorage](ctx, tx, listBuildingResourceStorageForBuildingSqlTemplate, building)
}
