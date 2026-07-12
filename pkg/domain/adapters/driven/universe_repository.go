package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

const (
	createUniverseQuery = `
INSERT INTO
	universe (id, name, created_at)
	VALUES ($1, $2, $3)`

	getUniverseQuery = `
SELECT
	id,
	name,
	created_at,
	version
FROM
	universe
WHERE
	id = $1`

	listResourceQuery = `
SELECT
	id,
	name,
	start_amount,
	start_production,
	start_storage,
	created_at
FROM
	resource
ORDER BY
	created_at,
	resource`

	listUniverseQuery = `
SELECT
	id,
	name,
	created_at,
	version
FROM
	universe
ORDER BY
	created_at,
	name`

	deleteUniverseQuery = `DELETE FROM universe WHERE id = $1`
)

type UniverseRepository struct {
	conn db.Connection
}

func NewUniverseRepository(conn db.Connection) *UniverseRepository {
	return &UniverseRepository{
		conn: conn,
	}
}

func (r *UniverseRepository) Create(ctx context.Context, universe models.Universe) error {
	_, err := r.conn.Exec(ctx, createUniverseQuery, universe.Id, universe.Name, universe.CreatedAt.UTC())
	return parseDbError(err)
}

func (r *UniverseRepository) Get(ctx context.Context, id uuid.UUID) (models.Universe, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return models.Universe{}, err
	}
	defer tx.Close(ctx)

	dbUniverse, err := db.QueryOneTx[mappers.DbUniverse](ctx, tx, getUniverseQuery, id)
	if err != nil {
		return models.Universe{}, parseDbError(err)
	}

	return loadUniverseDetails(ctx, tx, dbUniverse)
}

func (r *UniverseRepository) List(ctx context.Context) ([]models.Universe, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbUniverses, err := db.QueryAllTx[mappers.DbUniverse](ctx, tx, listUniverseQuery)
	if err != nil {
		return nil, err
	}

	universes := make([]models.Universe, 0, len(dbUniverses))
	for id := range dbUniverses {
		universe, err := loadUniverseDetails(ctx, tx, dbUniverses[id])
		if err != nil {
			return nil, err
		}

		universes = append(universes, universe)
	}

	return universes, nil
}

func (r *UniverseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.conn.Exec(ctx, deleteUniverseQuery, id)
	return err
}

func loadUniverseDetails(ctx context.Context, tx db.Transaction, dbUniverse mappers.DbUniverse) (models.Universe, error) {
	universe := dbUniverse.ToDomain()

	var err error
	universe.Resources, err = db.QueryAllTx[models.Resource](
		ctx,
		tx,
		listResourceQuery,
	)
	if err != nil {
		return universe, err
	}

	universe.Buildings, err = loadBuildings(ctx, tx)
	if err != nil {
		return universe, err
	}

	return universe, nil
}
