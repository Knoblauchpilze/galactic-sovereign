package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
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
	building = $1
`

func (r *buildingCostRepositoryImpl) ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingCost, error) {
	res := tx.Query(ctx, listBuildingCostForBuildingSqlTemplate, building)
	if err := res.Err(); err != nil {
		return []persistence.BuildingCost{}, err
	}

	var out []persistence.BuildingCost
	parser := func(rows db.Scannable) error {
		var cost persistence.BuildingCost
		err := rows.Scan(&cost.Building, &cost.Resource, &cost.Cost, &cost.Progress)
		if err != nil {
			return err
		}

		out = append(out, cost)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.BuildingCost{}, err
	}

	return out, nil
}
