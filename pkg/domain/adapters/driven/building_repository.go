package driven

import (
	"context"
	"fmt"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
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

	listBuildingCostForBuildingQuery = `
SELECT
	resource,
	cost,
	progress
FROM
	building_cost
WHERE
	building = $1`
)

type buildingRepositoryImpl struct {
	conn db.Connection
}

func NewBuildingRepository(conn db.Connection) driven.ForListingBuildings {
	return &buildingRepositoryImpl{
		conn: conn,
	}
}

func (r *buildingRepositoryImpl) List(ctx context.Context) ([]models.Building, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbBuildings, err := db.QueryAllTx[mappers.DbBuilding](ctx, tx, listBuildingQuery)
	if err != nil {
		return nil, err
	}

	buildings := make([]models.Building, 0, len(dbBuildings))
	for id := range dbBuildings {
		building := dbBuildings[id].ToDomain()
		building.BaseCosts, err = db.QueryAllTx[models.BuildingCost](ctx, tx, listBuildingCostForBuildingQuery, dbBuildings[id].Id)
		if err != nil {
			return nil, err
		}

		fmt.Printf("found %d costs for %s\n", len(building.BaseCosts), building.Id)

		buildings = append(buildings, building)
	}

	return buildings, nil
}
