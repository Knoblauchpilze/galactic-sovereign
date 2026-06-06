package drivenadapter

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenport "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
)

const (
	listResourceQuery = `
SELECT
	id,
	name,
	start_amount,
	start_production,
	start_storage,
	created_at
FROM
	resource`
)

type resourceRepositoryImpl struct {
	conn db.Connection
}

func NewResourceRepository(conn db.Connection) drivenport.ForListingResources {
	return &resourceRepositoryImpl{
		conn: conn,
	}
}

func (r *resourceRepositoryImpl) List(ctx context.Context) ([]models.Resource, error) {
	return db.QueryAll[models.Resource](ctx, r.conn, listResourceQuery)
}
