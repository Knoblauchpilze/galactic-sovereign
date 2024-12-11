package repositories

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionRepository interface {
	Create(ctx context.Context, tx db.Transaction, action persistence.BuildingAction) (persistence.BuildingAction, error)
	Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.BuildingAction, error)
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.BuildingAction, error)
	ListBeforeCompletionTime(ctx context.Context, tx db.Transaction, planet uuid.UUID, until time.Time) ([]persistence.BuildingAction, error)
	Delete(ctx context.Context, tx db.Transaction, action uuid.UUID) error
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
	DeleteForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) error
}

type buildingActionRepositoryImpl struct{}

func NewBuildingActionRepository() BuildingActionRepository {
	return &buildingActionRepositoryImpl{}
}

const createBuildingActionSqlTemplate = `
INSERT INTO
	building_action (id, planet, building, current_level, desired_level, created_at, completed_at)
	VALUES($1, $2, $3, $4, $5, $6, $7)
`

func (r *buildingActionRepositoryImpl) Create(ctx context.Context, tx db.Transaction, action persistence.BuildingAction) (persistence.BuildingAction, error) {
	_, err := tx.Exec(ctx, createBuildingActionSqlTemplate, action.Id, action.Planet, action.Building, action.CurrentLevel, action.DesiredLevel, action.CreatedAt, action.CompletedAt)
	return action, err
}

const getBuildingActionSqlTemplate = `
SELECT
	id,
	planet,
	building,
	current_level,
	desired_level,
	created_at,
	completed_at
FROM
	building_action
WHERE
	id = $1`

func (r *buildingActionRepositoryImpl) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.BuildingAction, error) {
	return db.QueryOneTx[persistence.BuildingAction](ctx, tx, getBuildingActionSqlTemplate, id)
}

const listBuildingActionForPlanetSqlTemplate = `
SELECT
	id,
	planet,
	building,
	current_level,
	desired_level,
	created_at,
	completed_at
FROM
	building_action
WHERE
	planet = $1`

func (r *buildingActionRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.BuildingAction, error) {
	return db.QueryAllTx[persistence.BuildingAction](ctx, tx, listBuildingActionForPlanetSqlTemplate, planet)
}

const listBuildingActionBeforeCompletionTimeSqlTemplate = `
SELECT
	id,
	planet,
	building,
	current_level,
	desired_level,
	created_at,
	completed_at
FROM
	building_action
WHERE
	completed_at <= $1
	AND planet = $2`

func (r *buildingActionRepositoryImpl) ListBeforeCompletionTime(ctx context.Context, tx db.Transaction, planet uuid.UUID, until time.Time) ([]persistence.BuildingAction, error) {
	return db.QueryAllTx[persistence.BuildingAction](ctx, tx, listBuildingActionBeforeCompletionTimeSqlTemplate, until, planet)
}

const deleteBuildingActionCostsSqlTemplate = `DELETE FROM building_action_cost WHERE action = $1`
const deleteBuildingActionResourceProductionSqlTemplate = `DELETE FROM building_action_resource_production WHERE action = $1`
const deleteBuildingActionSqlTemplate = `DELETE FROM building_action WHERE id = $1`

func (r *buildingActionRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, action uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionCostsSqlTemplate, action)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionResourceProductionSqlTemplate, action)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionSqlTemplate, action)
	return err
}

// https://stackoverflow.com/questions/21662726/delete-using-left-outer-join-in-postgres
const deleteBuildingActionCostForPlanetSqlTemplate = `
DELETE FROM
	building_action_cost AS bacd
USING
	building_action_cost AS bac
	LEFT JOIN building_action AS ba ON ba.id = bac.action
WHERE
	bacd.action = bac.action
	AND ba.planet = $1`
const deleteBuildingActionResourceProductionForPlanetSqlTemplate = `
DELETE FROM
	building_action_resource_production AS barpd
USING
	building_action_resource_production AS barp
	LEFT JOIN building_action AS ba ON ba.id = barp.action
WHERE
	barpd.action = barp.action
	AND ba.planet = $1`
const deleteBuildingActionForPlanetSqlTemplate = `DELETE FROM building_action WHERE planet = $1`

func (r *buildingActionRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionCostForPlanetSqlTemplate, planet)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionResourceProductionForPlanetSqlTemplate, planet)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionForPlanetSqlTemplate, planet)
	return err
}

const deleteBuildingActionCostForPlayerSqlTemplate = `
DELETE FROM
	building_action_cost AS bacd
USING
	building_action_cost AS bac
	LEFT JOIN building_action AS ba ON ba.id = bac.action
	LEFT JOIN planet AS p ON p.id = ba.planet
WHERE
	bacd.action = bac.action
	AND p.player = $1`
const deleteBuildingActionResourceProductionForPlayerSqlTemplate = `
DELETE FROM
	building_action_resource_production AS barpd
USING
	building_action_resource_production AS barp
	LEFT JOIN building_action AS ba ON ba.id = barp.action
	LEFT JOIN planet AS p ON p.id = ba.planet
WHERE
	barpd.action = barp.action
	AND p.player = $1`
const deleteBuildingActionForPlayerSqlTemplate = `
DELETE FROM
	building_action AS bad
USING
	building_action AS ba
	LEFT JOIN planet AS p ON p.id = ba.planet
WHERE
	bad.id = ba.id
	AND p.player = $1`

func (r *buildingActionRepositoryImpl) DeleteForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionCostForPlayerSqlTemplate, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionResourceProductionForPlayerSqlTemplate, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteBuildingActionForPlayerSqlTemplate, player)
	return err
}
