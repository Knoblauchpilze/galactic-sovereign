package repositories

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionCostRepository interface {
	Create(ctx context.Context, tx db.Transaction, cost persistence.BuildingActionCost) (persistence.BuildingActionCost, error)
	ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionCost, error)
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
