package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
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

const createPlanetResourceProductionSqlTemplate = `
INSERT INTO
	planet_resource_production (planet, building, resource, production, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5, $6)`

func (r *planetResourceProductionRepositoryImpl) Create(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error) {
	_, err := tx.Exec(ctx, createPlanetResourceProductionSqlTemplate, production.Planet, production.Building, production.Resource, production.Production, production.CreatedAt, production.CreatedAt)
	production.UpdatedAt = production.CreatedAt
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
	AND building = $2`

func (r *planetResourceProductionRepositoryImpl) GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building *uuid.UUID) (persistence.PlanetResourceProduction, error) {
	return db.QueryOneTx[persistence.PlanetResourceProduction](ctx, tx, listPlanetResourceProductionForPlanetAndBuildingSqlTemplate, planet, building)
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
	planet = $1`

func (r *planetResourceProductionRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceProduction, error) {
	return db.QueryAllTx[persistence.PlanetResourceProduction](ctx, tx, listPlanetResourceProductionForPlanetSqlTemplate, planet)
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
	AND version = $7`

const updatePlanetResourceProductionWithoutBuildingSqlTemplate = `
UPDATE
	planet_resource_production
SET
	production = $1,
	updated_at = $2,
	version = $3
WHERE
	planet = $4
	AND building IS NULL
	AND resource = $5
	AND version = $6`

func (r *planetResourceProductionRepositoryImpl) Update(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error) {
	version := production.Version + 1

	var affectedRows int64
	var err error

	if production.Building != nil {
		affectedRows, err = tx.Exec(
			ctx,
			updatePlanetResourceProductionSqlTemplate,
			production.Production,
			production.UpdatedAt,
			version,
			production.Planet,
			production.Building,
			production.Resource,
			production.Version,
		)
	} else {
		affectedRows, err = tx.Exec(
			ctx,
			updatePlanetResourceProductionWithoutBuildingSqlTemplate,
			production.Production,
			production.UpdatedAt,
			version,
			production.Planet,
			production.Resource,
			production.Version,
		)
	}

	if err != nil {
		return production, err
	}
	if affectedRows != 1 {
		return production, errors.NewCode(OptimisticLockException)
	}

	production.Version = version

	return production, nil
}

const deletePlanetResourceProductionSqlTemplate = `DELETE FROM planet_resource_production WHERE planet = $1`

func (r *planetResourceProductionRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetResourceProductionSqlTemplate, planet)
	return err
}
