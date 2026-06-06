package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
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

	listBuildingResourceProductionForBuildingQuery = `
SELECT
	resource,
	base,
	progress
FROM
	building_resource_production
WHERE
	building = $1`

	listBuildingResourceStorageForBuildingQuery = `
SELECT
	resource,
	base,
	scale,
	progress
FROM
	building_resource_storage
WHERE
	building = $1`
)

type buildingRepositoryImpl struct {
	conn db.Connection
}

func NewBuildingRepository(conn db.Connection) drivenports.ForListingBuildings {
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
		building, err := loadBuildingDetails(ctx, tx, dbBuildings[id])
		if err != nil {
			return nil, err
		}

		buildings = append(buildings, building)
	}

	return buildings, nil
}

func loadBuildingDetails(ctx context.Context, tx db.Transaction, dbBuilding mappers.DbBuilding) (models.Building, error) {
	building := dbBuilding.ToDomain()

	var err error
	building.Costs, err = db.QueryAllTx[models.BuildingCost](
		ctx,
		tx,
		listBuildingCostForBuildingQuery,
		dbBuilding.Id,
	)
	if err != nil {
		return building, err
	}

	building.Productions, err = db.QueryAllTx[models.BuildingResourceProduction](
		ctx,
		tx,
		listBuildingResourceProductionForBuildingQuery,
		dbBuilding.Id,
	)
	if err != nil {
		return building, err
	}

	building.Storages, err = db.QueryAllTx[models.BuildingResourceStorage](
		ctx,
		tx,
		listBuildingResourceStorageForBuildingQuery,
		dbBuilding.Id,
	)
	if err != nil {
		return building, err
	}

	return building, nil
}
