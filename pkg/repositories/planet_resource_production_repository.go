package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceProductionRepository interface {
	Create(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error)
	GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building *uuid.UUID) (persistence.PlanetResourceProduction, error)
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceProduction, error)
	Update(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error)
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
}

type planetResourceProductionRepositoryImpl struct{}

func NewPlanetResourceProductionRepository() PlanetResourceProductionRepository {
	return &planetResourceProductionRepositoryImpl{}
}

const createPlanetResourceProductionSqlTemplate = "INSERT INTO planet_resource_production (planet, building, resource, production, created_at) VALUES($1, $2, $3, $4, $5)"

func (r *planetResourceProductionRepositoryImpl) Create(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error) {
	_, err := tx.Exec(ctx, createPlanetResourceProductionSqlTemplate, production.Planet, production.Building, production.Resource, production.Production, production.CreatedAt)
	return production, err
}

const listPlanetResourceProductionForPlanetAndBuildingSqlTemplate = `
SELECT
	planet,
	building,
	resource,
	production,
	created_at,
	updated_at,
	version
FROM
	planet_resource_production
WHERE
	planet = $1
	AND building = $2
`

func (r *planetResourceProductionRepositoryImpl) GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building *uuid.UUID) (persistence.PlanetResourceProduction, error) {
	res := tx.Query(ctx, listPlanetResourceProductionForPlanetAndBuildingSqlTemplate, planet, building)
	if err := res.Err(); err != nil {
		return persistence.PlanetResourceProduction{}, err
	}

	var out persistence.PlanetResourceProduction
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Planet, &out.Building, &out.Resource, &out.Production, &out.CreatedAt, &out.UpdatedAt, &out.Version)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.PlanetResourceProduction{}, err
	}

	return out, nil
}

const listPlanetResourceProductionForPlanetSqlTemplate = `
SELECT
	planet,
	building,
	resource,
	production,
	created_at,
	updated_at,
	version
FROM
	planet_resource_production
WHERE
	planet = $1
`

func (r *planetResourceProductionRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceProduction, error) {
	res := tx.Query(ctx, listPlanetResourceProductionForPlanetSqlTemplate, planet)
	if err := res.Err(); err != nil {
		return []persistence.PlanetResourceProduction{}, err
	}

	var out []persistence.PlanetResourceProduction
	parser := func(rows db.Scannable) error {
		var production persistence.PlanetResourceProduction
		err := rows.Scan(&production.Planet, &production.Building, &production.Resource, &production.Production, &production.CreatedAt, &production.UpdatedAt, &production.Version)
		if err != nil {
			return err
		}

		out = append(out, production)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.PlanetResourceProduction{}, err
	}

	return out, nil
}

const updatePlanetResourceProductionSqlTemplate = `
UPDATE
	planet_resource_production
SET
	production = $1,
	updated_at = $2,
	version = $3
WHERE
	planet = $4
	AND building = $5
	AND resource = $6
	AND version = $7
`

func (r *planetResourceProductionRepositoryImpl) Update(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error) {
	version := production.Version + 1
	affected, err := tx.Exec(ctx, updatePlanetResourceProductionSqlTemplate,
		production.Production, production.UpdatedAt, version,
		production.Planet, production.Building, production.Resource, production.Version)
	if err != nil {
		return production, err
	}
	if affected == 0 {
		return production, errors.NewCode(db.OptimisticLockException)
	} else if affected != 1 {
		return production, errors.NewCode(db.MoreThanOneMatchingSqlRows)
	}

	production.Version = version

	return production, nil
}

const deletePlanetResourceProductionSqlTemplate = "DELETE FROM planet_resource_production WHERE planet = $1"

func (r *planetResourceProductionRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetResourceProductionSqlTemplate, planet)
	return err
}
