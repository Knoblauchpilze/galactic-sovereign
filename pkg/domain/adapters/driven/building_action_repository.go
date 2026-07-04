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
	upsertBuildingActionQuery = `
INSERT INTO
	building_action (id, planet, building, desired_level, created_at, completed_at)
	VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE
SET
	completed_at = excluded.completed_at`

	// For this query and the following ones, voluntarily doing nothing on conflict.
	// This allows the action dependencies to behave as a single aggregate for which
	// only the completion time can be updated.
	upsertBuildingActionCostQuery = `
INSERT INTO
	building_action_cost (action, resource, amount)
	VALUES ($1, $2, $3)
ON CONFLICT (action, resource) DO NOTHING`

	upsertBuildingActionResourceStorageQuery = `
INSERT INTO
	building_action_resource_storage (action, resource, storage)
	VALUES ($1, $2, $3)
ON CONFLICT (action, resource) DO NOTHING`

	upsertBuildingActionResourceProductionQuery = `
INSERT INTO
	building_action_resource_production (action, resource, production)
	VALUES ($1, $2, $3)
ON CONFLICT (action, resource) DO NOTHING`

	getBuildingActionQuery = `
SELECT
	id,
	building,
	desired_level,
	created_at,
	completed_at
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

// TODO: Deprecated, this should be removed when the planet mutator is ready
func NewBuildingActionRepository(conn db.Connection) drivenports.ForManagingBuildingActions {
	return &buildingActionRepositoryImpl{
		conn: conn,
	}
}

func (r *buildingActionRepositoryImpl) Create(
	ctx context.Context,
	planet models.Planet,
) error {
	if planet.BuildingAction == nil {
		return nil
	}

	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	err = upsertBuildingActionWithDetails(ctx, tx, planet.Id, *planet.BuildingAction)
	if err != nil {
		return parseDbError(err)
	}

	err = updatePlanetDetails(ctx, tx, planet, planet.Version-1)
	if err != nil {
		return parseDbError(err)
	}

	return nil
}

func (r *buildingActionRepositoryImpl) Delete(
	ctx context.Context,
	planet models.Planet,
) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	err = deleteBuildingActionAndDetailsForPlanet(ctx, tx, planet.Id)
	if err != nil {
		return err
	}

	err = updatePlanetDetails(ctx, tx, planet, planet.Version-1)
	if err != nil {
		return parseDbError(err)
	}

	return nil
}

func upsertBuildingActionWithDetails(
	ctx context.Context,
	tx db.Transaction,
	planet uuid.UUID,
	action models.BuildingAction,
) error {
	_, err := tx.Exec(
		ctx,
		upsertBuildingActionQuery,
		action.Id,
		planet,
		action.Building,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
	)
	if err != nil {
		return err
	}

	for _, c := range action.Costs {
		_, err = tx.Exec(
			ctx,
			upsertBuildingActionCostQuery,
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
			upsertBuildingActionResourceStorageQuery,
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
			upsertBuildingActionResourceProductionQuery,
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

func loadBuildingActionAndDetails(
	ctx context.Context,
	tx db.Transaction,
	id uuid.UUID,
) (models.BuildingAction, error) {
	dbAction, err := db.QueryOneTx[mappers.DbBuildingAction](
		ctx,
		tx,
		getBuildingActionQuery,
		id,
	)
	if err != nil {
		return models.BuildingAction{}, parseDbError(err)
	}

	action := dbAction.ToDomain()

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

func deleteBuildingActionAndDetailsForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
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
