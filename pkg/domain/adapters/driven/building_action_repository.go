package drivenadapters

import (
	"context"
	"time"

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

	listBuildingActionForPlanetQuery = `
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
	planet = $1`

	listBuildingActionBeforeCompletionTimeQuery = `
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
	planet = $1
	AND completed_at <= $2`

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

	deleteBuildingActionResourceProductionForPlayerQuery = `
DELETE FROM
	building_action_resource_production AS barpd
USING
	building_action_resource_production AS barp
	LEFT JOIN building_action AS ba ON ba.id = barp.action
	LEFT JOIN planet AS p ON p.id = ba.planet
WHERE
	barpd.action = barp.action
	AND p.player = $1`
	deleteBuildingActionResourceStorageForPlayerQuery = `
DELETE FROM
	building_action_resource_storage AS barsd
USING
	building_action_resource_storage AS bars
	LEFT JOIN building_action AS ba ON ba.id = bars.action
	LEFT JOIN planet AS p ON p.id = ba.planet
WHERE
	barsd.action = bars.action
	AND p.player = $1`
	deleteBuildingActionCostForPlayerQuery = `
DELETE FROM
	building_action_cost AS bacd
USING
	building_action_cost AS bac
	LEFT JOIN building_action AS ba ON ba.id = bac.action
	LEFT JOIN planet AS p ON p.id = ba.planet
WHERE
	bacd.action = bac.action
	AND p.player = $1`
	deleteBuildingActionForPlayerQuery = `
DELETE FROM
	building_action AS bad
USING
	building_action AS ba
	LEFT JOIN planet AS p ON p.id = ba.planet
WHERE
	bad.id = ba.id
	AND p.player = $1`
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
	_, err := r.conn.Exec(
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
	return err
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
		return models.BuildingAction{}, err
	}

	return loadBuildingActionDetails(ctx, tx, dbAction)
}

func (r *buildingActionRepositoryImpl) ListForPlanet(
	ctx context.Context,
	planet uuid.UUID,
) ([]models.BuildingAction, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbActions, err := db.QueryAllTx[mappers.DbBuildingAction](
		ctx,
		tx,
		listBuildingActionForPlanetQuery,
		planet,
	)
	if err != nil {
		return nil, err
	}

	actions := make([]models.BuildingAction, 0, len(dbActions))
	for id := range dbActions {
		action, err := loadBuildingActionDetails(ctx, tx, dbActions[id])
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func (r *buildingActionRepositoryImpl) ListBeforeCompletionTime(
	ctx context.Context,
	planet uuid.UUID,
	until time.Time,
) ([]models.BuildingAction, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbActions, err := db.QueryAllTx[mappers.DbBuildingAction](
		ctx,
		tx,
		listBuildingActionBeforeCompletionTimeQuery,
		planet,
		until,
	)
	if err != nil {
		return nil, err
	}

	actions := make([]models.BuildingAction, 0, len(dbActions))
	for id := range dbActions {
		action, err := loadBuildingActionDetails(ctx, tx, dbActions[id])
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
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

func (r *buildingActionRepositoryImpl) DeleteForPlanet(ctx context.Context, planet uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deleteBuildingActionResourceProductionForPlanetQuery, planet)
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

func (r *buildingActionRepositoryImpl) DeleteForPlayer(ctx context.Context, player uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deleteBuildingActionResourceProductionForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionResourceStorageForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionCostForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionForPlayerQuery, player)
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
