package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
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
	building = $1
`

func (r *buildingResourceProductionRepositoryImpl) ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingResourceProduction, error) {
	res := tx.Query(ctx, listBuildingResourceProductionForBuildingSqlTemplate, building)
	if err := res.Err(); err != nil {
		return []persistence.BuildingResourceProduction{}, err
	}

	var out []persistence.BuildingResourceProduction
	parser := func(rows db.Scannable) error {
		var cost persistence.BuildingResourceProduction
		err := rows.Scan(&cost.Building, &cost.Resource, &cost.Base, &cost.Progress)
		if err != nil {
			return err
		}

		out = append(out, cost)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.BuildingResourceProduction{}, err
	}

	return out, nil
}
