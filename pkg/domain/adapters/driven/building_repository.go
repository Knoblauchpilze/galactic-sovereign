package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

const (
	getBuildingQuery = `
SELECT
	id,
	name,
	created_at
FROM
	building
WHERE
	id = $1`

	listBuildingQuery = `
SELECT
	id,
	name,
	created_at
FROM
	building
ORDER BY
	created_at,
	name`

	listBuildingCostForBuildingQuery = `
SELECT
	bc.resource,
	bc.cost,
	bc.progress,
	r.build_time_hours_per_unit
FROM
	building_cost AS bc
	INNER JOIN resource AS r ON r.id = bc.resource
WHERE
	bc.building = $1`

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

type BuildingRepository struct {
	conn db.Connection
}

func NewBuildingRepository(conn db.Connection) *BuildingRepository {
	return &BuildingRepository{
		conn: conn,
	}
}

func (r *BuildingRepository) Get(ctx context.Context, id uuid.UUID) (models.Building, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return models.Building{}, err
	}
	defer tx.Close(ctx)

	dbBuilding, err := db.QueryOneTx[mappers.DbBuilding](ctx, tx, getBuildingQuery, id)
	if err != nil {
		return models.Building{}, parseDbError(err)
	}

	return loadBuildingDetails(ctx, tx, dbBuilding)
}

func loadBuildings(ctx context.Context, tx db.Transaction) ([]models.Building, error) {
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
