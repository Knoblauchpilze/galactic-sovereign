package repositories

import (
	"context"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
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
	if err != nil && duplicatedKeySqlErrorRegexp.MatchString(err.Error()) {
		return production, errors.NewCode(db.DuplicatedKeySqlKey)
	}

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
	res := tx.Query(ctx, listBuildingActionResourceProductionForActionSqlTemplate, action)
	if err := res.Err(); err != nil {
		return []persistence.BuildingActionResourceProduction{}, err
	}

	var out []persistence.BuildingActionResourceProduction
	parser := func(rows db.Scannable) error {
		var production persistence.BuildingActionResourceProduction
		err := rows.Scan(&production.Action, &production.Resource, &production.Production)
		if err != nil {
			return err
		}

		out = append(out, production)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.BuildingActionResourceProduction{}, err
	}

	return out, nil
}

const deleteBuildingActionResourceProductionForActionSqlTemplate = "DELETE FROM building_action_resource_production WHERE action = $1"

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
