package repositories

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceRepository interface {
	Create(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error)
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResource, error)
	Update(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error)
}

type planetResourceRepositoryImpl struct{}

func NewPlanetResourceRepository() PlanetResourceRepository {
	return &planetResourceRepositoryImpl{}
}

const createPlanetResourceSqlTemplate = `
INSERT INTO
	planet_resource (planet, resource, amount, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5)`

func (r *planetResourceRepositoryImpl) Create(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error) {
	_, err := tx.Exec(ctx, createPlanetResourceSqlTemplate, resource.Planet, resource.Resource, resource.Amount, resource.CreatedAt, resource.CreatedAt)
	resource.UpdatedAt = resource.CreatedAt
	return resource, err
}

const listPlanetResourceForPlanetSqlTemplate = `
SELECT
	planet,
	resource,
	amount,
	created_at,
	updated_at,
	version
FROM
	planet_resource
WHERE
	planet = $1`

func (r *planetResourceRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResource, error) {
	return db.QueryAllTx[persistence.PlanetResource](ctx, tx, listPlanetResourceForPlanetSqlTemplate, planet)
}

const updatePlanetResourceSqlTemplate = `
UPDATE
	planet_resource
SET
	amount = $1,
	updated_at = $2,
	version = $3
WHERE
	planet = $4
	AND resource = $5
	AND version = $6`

func (r *planetResourceRepositoryImpl) Update(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error) {
	version := resource.Version + 1

	affectedRows, err := tx.Exec(ctx, updatePlanetResourceSqlTemplate, resource.Amount, resource.UpdatedAt, version, resource.Planet, resource.Resource, resource.Version)
	if err != nil {
		return resource, err
	}
	if affectedRows != 1 {
		return resource, errors.NewCode(OptimisticLockException)
	}

	resource.Version = version

	return resource, nil
}
