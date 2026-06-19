package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

const (
	createBuildingActionQuery = `
INSERT INTO
	building_action (id, planet, building, current_level, desired_level, created_at, completed_at, version)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	createBuildingActionCostQuery = `
INSERT INTO
	building_action_cost (action, resource, amount)
	VALUES($1, $2, $3)`

	createBuildingActionResourceStorageQuery = `
INSERT INTO
	building_action_resource_storage (action, resource, storage)
	VALUES($1, $2, $3)`

	createBuildingActionResourceProductionQuery = `
INSERT INTO
	building_action_resource_production (action, resource, production)
	VALUES($1, $2, $3)`

	getBuildingActionQuery = `
SELECT
	id,
	planet,
	building,
	current_level,
	desired_level,
	created_at,
	completed_at,
	version
FROM
	building_action
WHERE
	id = $1`

	listBuildingActionCostForActionQuery = `
SELECT
	resource,
	amount
FROM
	building_action_cost
WHERE
	action = $1`

	listBuildingActionResourceStorageForActionQuery = `
SELECT
	resource,
	storage
FROM
	building_action_resource_storage
WHERE
	action = $1`

	listBuildingActionResourceProductionForActionQuery = `
SELECT
	resource,
	production
FROM
	building_action_resource_production
WHERE
	action = $1`

	deleteBuildingActionResourceProductionQuery = `DELETE FROM building_action_resource_production WHERE action = $1`
	deleteBuildingActionResourceStorageQuery    = `DELETE FROM building_action_resource_storage WHERE action = $1`
	deleteBuildingActionCostsQuery              = `DELETE FROM building_action_cost WHERE action = $1`
	deleteBuildingActionQuery                   = `DELETE FROM building_action WHERE id = $1`

	deleteBuildingActionResourceProductionForPlanetQuery = `
DELETE FROM
	building_action_resource_production AS barpd
USING
	building_action_resource_production AS barp
	LEFT JOIN building_action AS ba ON ba.id = barp.action
WHERE
	barpd.action = barp.action
	AND ba.planet = $1`
	deleteBuildingActionResourceStorageForPlanetQuery = `
DELETE FROM
	building_action_resource_storage AS barsd
USING
	building_action_resource_storage AS bars
	LEFT JOIN building_action AS ba ON ba.id = bars.action
WHERE
	barsd.action = bars.action
	AND ba.planet = $1`
	// https://stackoverflow.com/questions/21662726/delete-using-left-outer-join-in-postgres
	deleteBuildingActionCostForPlanetQuery = `
DELETE FROM
	building_action_cost AS bacd
USING
	building_action_cost AS bac
	LEFT JOIN building_action AS ba ON ba.id = bac.action
WHERE
	bacd.action = bac.action
	AND ba.planet = $1`
	deleteBuildingActionForPlanetQuery = `DELETE FROM building_action WHERE planet = $1`
)

type buildingActionRepositoryImpl struct {
	conn db.Connection
}

func NewBuildingActionRepository(conn db.Connection) drivenports.ForManagingBuildingActions {
	return &buildingActionRepositoryImpl{
		conn: conn,
	}
}

func (r *buildingActionRepositoryImpl) Create(
	ctx context.Context,
	action models.BuildingAction,
) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(
		ctx,
		createBuildingActionQuery,
		action.Id,
		action.Planet,
		action.Building,
		action.CurrentLevel,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
		action.Version,
	)
	if err != nil {
		return parseDbError(err)
	}

	for _, c := range action.Costs {
		_, err = tx.Exec(
			ctx,
			createBuildingActionCostQuery,
			action.Id,
			c.Resource,
			c.Amount,
		)
		if err != nil {
			return err
		}
	}

	for _, s := range action.Storages {
		_, err = tx.Exec(
			ctx,
			createBuildingActionResourceStorageQuery,
			action.Id,
			s.Resource,
			s.Storage,
		)
		if err != nil {
			return err
		}
	}

	for _, p := range action.Productions {
		_, err = tx.Exec(
			ctx,
			createBuildingActionResourceProductionQuery,
			action.Id,
			p.Resource,
			p.Production,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *buildingActionRepositoryImpl) Get(
	ctx context.Context,
	id uuid.UUID,
) (models.BuildingAction, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return models.BuildingAction{}, err
	}
	defer tx.Close(ctx)

	dbAction, err := db.QueryOneTx[mappers.DbBuildingAction](
		ctx,
		tx,
		getBuildingActionQuery,
		id,
	)
	if err != nil {
		return models.BuildingAction{}, parseDbError(err)
	}

	return loadBuildingActionDetails(ctx, tx, dbAction)
}

func (r *buildingActionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deleteBuildingActionResourceProductionQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionResourceStorageQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionCostsQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func loadBuildingActionDetails(
	ctx context.Context,
	tx db.Transaction,
	dbAction mappers.DbBuildingAction,
) (models.BuildingAction, error) {
	action := dbAction.ToDomain()

	var err error
	action.Costs, err = db.QueryAllTx[models.BuildingActionCost](
		ctx,
		tx,
		listBuildingActionCostForActionQuery,
		dbAction.Id,
	)
	if err != nil {
		return action, err
	}

	action.Storages, err = db.QueryAllTx[models.BuildingActionResourceStorage](
		ctx,
		tx,
		listBuildingActionResourceStorageForActionQuery,
		dbAction.Id,
	)
	if err != nil {
		return action, err
	}

	action.Productions, err = db.QueryAllTx[models.BuildingActionResourceProduction](
		ctx,
		tx,
		listBuildingActionResourceProductionForActionQuery,
		dbAction.Id,
	)
	if err != nil {
		return action, err
	}

	return action, nil
}

func deleteBuildingActionDetailsForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionResourceProductionForPlanetQuery, planet)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionResourceStorageForPlanetQuery, planet)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionCostForPlanetQuery, planet)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionForPlanetQuery, planet)
	if err != nil {
		return err
	}

	return nil
}
