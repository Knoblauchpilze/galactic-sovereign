package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
)

type resourceRepositoryImpl struct {
	conn db.Connection
}

func NewResourceRepository(conn db.Connection) drivenports.ForListingResources {
	return &resourceRepositoryImpl{
		conn: conn,
	}
}

func (r *resourceRepositoryImpl) List(ctx context.Context) ([]models.Resource, error) {
	return db.QueryAll[models.Resource](ctx, r.conn, listResourceQuery)
}
