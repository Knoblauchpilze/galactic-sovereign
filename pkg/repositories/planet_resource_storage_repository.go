package repositories

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceStorageRepository interface {
	Create(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error)
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceStorage, error)
	Update(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error)
}

type planetResourceStorageRepositoryImpl struct{}

func NewPlanetResourceStorageRepository() PlanetResourceStorageRepository {
	return &planetResourceStorageRepositoryImpl{}
}

const createPlanetResourceStorageSqlTemplate = `
INSERT INTO
	planet_resource_storage (planet, resource, storage, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5)`

func (r *planetResourceStorageRepositoryImpl) Create(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error) {
	_, err := tx.Exec(ctx, createPlanetResourceStorageSqlTemplate, storage.Planet, storage.Resource, storage.Storage, storage.CreatedAt, storage.CreatedAt)
	storage.UpdatedAt = storage.CreatedAt
	return storage, err
}

const listPlanetResourceStorageForPlanetSqlTemplate = `
SELECT
	planet,
	resource,
	storage,
	created_at,
	updated_at,
	version
FROM
	planet_resource_storage
WHERE
	planet = $1`

func (r *planetResourceStorageRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceStorage, error) {
	return db.QueryAllTx[persistence.PlanetResourceStorage](ctx, tx, listPlanetResourceStorageForPlanetSqlTemplate, planet)
}

const updatePlanetResourceStorageSqlTemplate = `
UPDATE
	planet_resource_storage
SET
	storage = $1,
	updated_at = $2,
	version = $3
WHERE
	planet = $4
	AND resource = $5
	AND version = $6`

func (r *planetResourceStorageRepositoryImpl) Update(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error) {
	version := storage.Version + 1
	affectedRows, err := tx.Exec(
		ctx,
		updatePlanetResourceStorageSqlTemplate,
		storage.Storage,
		storage.UpdatedAt,
		version,
		storage.Planet,
		storage.Resource,
		storage.Version,
	)
	if err != nil {
		return storage, err
	}
	if affectedRows != 1 {
		return storage, errors.NewCode(OptimisticLockException)
	}

	storage.Version = version

	return storage, nil
}
