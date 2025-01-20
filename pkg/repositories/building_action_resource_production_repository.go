package repositories

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionResourceProductionRepository interface {
	Create(ctx context.Context, tx db.Transaction, production persistence.BuildingActionResourceProduction) (persistence.BuildingActionResourceProduction, error)
	ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionResourceProduction, error)
}

type buildingActionResourceProductionRepositoryImpl struct{}

func NewBuildingActionResourceProductionRepository() BuildingActionResourceProductionRepository {
	return &buildingActionResourceProductionRepositoryImpl{}
}

const createBuildingActionResourceProductionSqlTemplate = `
INSERT INTO
	building_action_resource_production (action, resource, production)
	VALUES($1, $2, $3)`

func (r *buildingActionResourceProductionRepositoryImpl) Create(ctx context.Context, tx db.Transaction, production persistence.BuildingActionResourceProduction) (persistence.BuildingActionResourceProduction, error) {
	_, err := tx.Exec(ctx, createBuildingActionResourceProductionSqlTemplate, production.Action, production.Resource, production.Production)
	return production, err
}

const listBuildingActionResourceProductionForActionSqlTemplate = `
SELECT
	action,
	resource,
	production
FROM
	building_action_resource_production
WHERE
	action = $1`

func (r *buildingActionResourceProductionRepositoryImpl) ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionResourceProduction, error) {
	return db.QueryAllTx[persistence.BuildingActionResourceProduction](ctx, tx, listBuildingActionResourceProductionForActionSqlTemplate, action)
}
