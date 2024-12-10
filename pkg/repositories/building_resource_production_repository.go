package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingResourceProductionRepository interface {
	ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingResourceProduction, error)
}

type buildingResourceProductionRepositoryImpl struct{}

func NewBuildingResourceProductionRepository() BuildingResourceProductionRepository {
	return &buildingResourceProductionRepositoryImpl{}
}

const listBuildingResourceProductionForBuildingSqlTemplate = `
SELECT
	building,
	resource,
	base,
	progress
FROM
	building_resource_production
WHERE
	building = $1`

func (r *buildingResourceProductionRepositoryImpl) ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingResourceProduction, error) {
	return db.QueryAllTx[persistence.BuildingResourceProduction](ctx, tx, listBuildingResourceProductionForBuildingSqlTemplate, building)
}
