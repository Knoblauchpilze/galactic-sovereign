package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionResourceProductionRepository interface {
	Create(ctx context.Context, tx db.Transaction, production persistence.BuildingActionResourceProduction) (persistence.BuildingActionResourceProduction, error)
	ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionResourceProduction, error)
	DeleteForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) error
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
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

const deleteBuildingActionResourceProductionForActionSqlTemplate = `DELETE FROM building_action_resource_production WHERE action = $1`

func (r *buildingActionResourceProductionRepositoryImpl) DeleteForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionResourceProductionForActionSqlTemplate, action)
	return err
}

const deleteBuildingActionResourceProductionForPlanetSqlTemplate = `
DELETE FROM
	building_action_resource_production
USING
	building_action_resource_production AS barp
	LEFT JOIN building_action AS ba ON ba.id = barp.action
WHERE
	building_action_resource_production.action = barp.action
	AND ba.planet = $1`

func (r *buildingActionResourceProductionRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionResourceProductionForPlanetSqlTemplate, planet)
	return err
}
