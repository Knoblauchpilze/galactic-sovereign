package repositories

import (
	"context"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
)

type BuildingRepository interface {
	List(ctx context.Context, tx db.Transaction) ([]persistence.Building, error)
}

type buildingRepositoryImpl struct{}

func NewBuildingRepository() BuildingRepository {
	return &buildingRepositoryImpl{}
}

const listBuildingSqlTemplate = "SELECT id, name, created_at, updated_at FROM building"

func (r *buildingRepositoryImpl) List(ctx context.Context, tx db.Transaction) ([]persistence.Building, error) {
	res := tx.Query(ctx, listBuildingSqlTemplate)
	if err := res.Err(); err != nil {
		return []persistence.Building{}, err
	}

	var out []persistence.Building
	parser := func(rows db.Scannable) error {
		var building persistence.Building
		err := rows.Scan(&building.Id, &building.Name, &building.CreatedAt, &building.UpdatedAt)
		if err != nil {
			return err
		}

		out = append(out, building)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.Building{}, err
	}

	return out, nil
}
