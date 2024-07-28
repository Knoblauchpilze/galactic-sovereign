package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceRepository interface {
	Create(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error)
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResource, error)
	Update(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error)
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
}

type planetResourceRepositoryImpl struct{}

func NewPlanetResourceRepository() PlanetResourceRepository {
	return &planetResourceRepositoryImpl{}
}

const createPlanetResourceSqlTemplate = "INSERT INTO planet_resource (planet, resource, amount, created_at) VALUES($1, $2, $3, $4)"

func (r *planetResourceRepositoryImpl) Create(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error) {
	_, err := tx.Exec(ctx, createPlanetResourceSqlTemplate, resource.Planet, resource.Resource, resource.Amount, resource.CreatedAt)
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
	planet = $1
`

func (r *planetResourceRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResource, error) {
	res := tx.Query(ctx, listPlanetResourceForPlanetSqlTemplate, planet)
	if err := res.Err(); err != nil {
		return []persistence.PlanetResource{}, err
	}

	var out []persistence.PlanetResource
	parser := func(rows db.Scannable) error {
		var resource persistence.PlanetResource
		err := rows.Scan(&resource.Planet, &resource.Resource, &resource.Amount, &resource.CreatedAt, &resource.UpdatedAt, &resource.Version)
		if err != nil {
			return err
		}

		out = append(out, resource)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.PlanetResource{}, err
	}

	return out, nil
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
	AND version = $6
`

func (r *planetResourceRepositoryImpl) Update(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error) {
	version := resource.Version + 1
	affected, err := tx.Exec(ctx, updatePlanetResourceSqlTemplate, resource.Amount, resource.UpdatedAt, version, resource.Planet, resource.Resource, resource.Version)
	if err != nil {
		return resource, err
	}
	if affected == 0 {
		return resource, errors.NewCode(db.OptimisticLockException)
	} else if affected != 1 {
		return resource, errors.NewCode(db.MoreThanOneMatchingSqlRows)
	}

	resource.Version = version

	return resource, nil
}

const deletePlanetResourceSqlTemplate = "DELETE FROM planet_resource WHERE planet = $1"

func (r *planetResourceRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetResourceSqlTemplate, planet)
	return err
}
