package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceStorageRepository interface {
	Create(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error)
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceStorage, error)
	Update(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error)
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
}

type planetResourceStorageRepositoryImpl struct{}

func NewPlanetResourceStorageRepository() PlanetResourceStorageRepository {
	return &planetResourceStorageRepositoryImpl{}
}

const createPlanetResourceStorageSqlTemplate = "INSERT INTO planet_resource_storage (planet, resource, storage, created_at) VALUES($1, $2, $3, $4)"

func (r *planetResourceStorageRepositoryImpl) Create(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error) {
	_, err := tx.Exec(ctx, createPlanetResourceStorageSqlTemplate, storage.Planet, storage.Resource, storage.Storage, storage.CreatedAt)
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
	planet = $1
`

func (r *planetResourceStorageRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceStorage, error) {
	res := tx.Query(ctx, listPlanetResourceStorageForPlanetSqlTemplate, planet)
	if err := res.Err(); err != nil {
		return []persistence.PlanetResourceStorage{}, err
	}

	var out []persistence.PlanetResourceStorage
	parser := func(rows db.Scannable) error {
		var storage persistence.PlanetResourceStorage
		err := rows.Scan(&storage.Planet, &storage.Resource, &storage.Storage, &storage.CreatedAt, &storage.UpdatedAt, &storage.Version)
		if err != nil {
			return err
		}

		out = append(out, storage)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.PlanetResourceStorage{}, err
	}

	return out, nil
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
	AND version = $6
`

func (r *planetResourceStorageRepositoryImpl) Update(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error) {
	version := storage.Version + 1
	affected, err := tx.Exec(ctx, updatePlanetResourceStorageSqlTemplate,
		storage.Storage, storage.UpdatedAt, version,
		storage.Planet, storage.Resource, storage.Version)
	if err != nil {
		return storage, err
	}
	if affected == 0 {
		return storage, errors.NewCode(db.OptimisticLockException)
	} else if affected != 1 {
		return storage, errors.NewCode(db.MoreThanOneMatchingSqlRows)
	}

	storage.Version = version

	return storage, nil
}

const deletePlanetResourceStorageSqlTemplate = "DELETE FROM planet_resource_storage WHERE planet = $1"

func (r *planetResourceStorageRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetResourceStorageSqlTemplate, planet)
	return err
}
