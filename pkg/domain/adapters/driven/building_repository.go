package driven

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
)

const (
	listBuildingQuery = `
SELECT
	id,
	name,
	created_at
FROM
	building`
)

type buildingRepositoryImpl struct {
	conn db.Connection
}

func NewBuildingRepository(conn db.Connection) driven.ForManagingBuildings {
	return &buildingRepositoryImpl{
		conn: conn,
	}
}

func (r *buildingRepositoryImpl) List(ctx context.Context) ([]models.Building, error) {
	return db.QueryAll[models.Building](ctx, r.conn, listBuildingQuery)
}
