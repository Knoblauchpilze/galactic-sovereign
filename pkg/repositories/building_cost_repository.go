package repositories

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingCostRepository interface {
	ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingCost, error)
}

type buildingCostRepositoryImpl struct{}

func NewBuildingCostRepository() BuildingCostRepository {
	return &buildingCostRepositoryImpl{}
}

const listBuildingCostForBuildingSqlTemplate = `
SELECT
	building,
	resource,
	cost,
	progress
FROM
	building_cost
WHERE
	building = $1`

func (r *buildingCostRepositoryImpl) ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingCost, error) {
	return db.QueryAllTx[persistence.BuildingCost](ctx, tx, listBuildingCostForBuildingSqlTemplate, building)
}
