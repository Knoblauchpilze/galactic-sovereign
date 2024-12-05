package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionCostRepository interface {
	Create(ctx context.Context, tx db.Transaction, cost persistence.BuildingActionCost) (persistence.BuildingActionCost, error)
	ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionCost, error)
	DeleteForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) error
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
}

type buildingActionCostRepositoryImpl struct{}

func NewBuildingActionCostRepository() BuildingActionCostRepository {
	return &buildingActionCostRepositoryImpl{}
}

const createBuildingActionCostSqlTemplate = `
INSERT INTO
	building_action_cost (action, resource, amount)
	VALUES($1, $2, $3)`

func (r *buildingActionCostRepositoryImpl) Create(ctx context.Context, tx db.Transaction, cost persistence.BuildingActionCost) (persistence.BuildingActionCost, error) {
	_, err := tx.Exec(ctx, createBuildingActionCostSqlTemplate, cost.Action, cost.Resource, cost.Amount)
	return cost, err
}

const listBuildingActionCostForActionSqlTemplate = `
SELECT
	action,
	resource,
	amount
FROM
	building_action_cost
WHERE
	action = $1`

func (r *buildingActionCostRepositoryImpl) ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionCost, error) {
	return db.QueryAllTx[persistence.BuildingActionCost](ctx, tx, listBuildingActionCostForActionSqlTemplate, action)
}

const deleteBuildingActionCostForActionSqlTemplate = `DELETE FROM building_action_cost WHERE action = $1`

func (r *buildingActionCostRepositoryImpl) DeleteForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionCostForActionSqlTemplate, action)
	return err
}

// https://stackoverflow.com/questions/21662726/delete-using-left-outer-join-in-postgres
const deleteBuildingActionCostForPlanetSqlTemplate = `
DELETE FROM
	building_action_cost
USING
	building_action_cost AS bac
	LEFT JOIN building_action AS ba ON ba.id = bac.action
WHERE
	building_action_cost.action = bac.action
	AND ba.planet = $1`

func (r *buildingActionCostRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionCostForPlanetSqlTemplate, planet)
	return err
}
