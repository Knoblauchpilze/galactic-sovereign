package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
)

type BuildingRepository interface {
	List(ctx context.Context, tx db.Transaction) ([]persistence.Building, error)
}

type buildingRepositoryImpl struct{}

func NewBuildingRepository() BuildingRepository {
	return &buildingRepositoryImpl{}
}

const listBuildingSqlTemplate = `
SELECT
	id,
	name,
	created_at,
	updated_at
FROM
	building`

func (r *buildingRepositoryImpl) List(ctx context.Context, tx db.Transaction) ([]persistence.Building, error) {
	return db.QueryAllTx[persistence.Building](ctx, tx, listBuildingSqlTemplate)
}
