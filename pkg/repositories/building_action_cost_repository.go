package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionCostRepository interface {
	Create(ctx context.Context, tx db.Transaction, cost persistence.BuildingActionCost) (persistence.BuildingActionCost, error)
	ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionCost, error)
	DeleteForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) error
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
	if err != nil && duplicatedKeySqlErrorRegexp.MatchString(err.Error()) {
		return cost, errors.NewCode(db.DuplicatedKeySqlKey)
	}

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
	res := tx.Query(ctx, listBuildingActionCostForActionSqlTemplate, action)
	if err := res.Err(); err != nil {
		return []persistence.BuildingActionCost{}, err
	}

	var out []persistence.BuildingActionCost
	parser := func(rows db.Scannable) error {
		var cost persistence.BuildingActionCost
		err := rows.Scan(&cost.Action, &cost.Resource, &cost.Amount)
		if err != nil {
			return err
		}

		out = append(out, cost)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.BuildingActionCost{}, err
	}

	return out, nil
}

const deleteBuildingActionCostForActionSqlTemplate = "DELETE FROM building_action_cost WHERE action = $1"

func (r *buildingActionCostRepositoryImpl) DeleteForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionCostForActionSqlTemplate, action)
	return err
}
