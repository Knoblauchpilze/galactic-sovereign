package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
)

type ResourceRepository interface {
	List(ctx context.Context, tx db.Transaction) ([]persistence.Resource, error)
}

type resourceRepositoryImpl struct{}

func NewResourceRepository() ResourceRepository {
	return &resourceRepositoryImpl{}
}

const listResourceSqlTemplate = `
SELECT
	id,
	name,
	start_amount,
	start_production,
	start_storage,
	created_at,
	updated_at
FROM
	resource`

func (r *resourceRepositoryImpl) List(ctx context.Context, tx db.Transaction) ([]persistence.Resource, error) {
	return db.QueryAllTx[persistence.Resource](ctx, tx, listResourceSqlTemplate)
}
