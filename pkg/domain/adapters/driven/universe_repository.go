package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

const (
	createUniverseQuery = `
INSERT INTO
	universe (id, name, created_at)
	VALUES($1, $2, $3)`

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

type universeRepositoryImpl struct {
	conn db.Connection
}

func NewUniverseRepository(conn db.Connection) drivenports.ForManagingUniverses {
	return &universeRepositoryImpl{
		conn: conn,
	}
}

func (r *universeRepositoryImpl) Create(ctx context.Context, universe models.Universe) error {
	_, err := r.conn.Exec(ctx, createUniverseQuery, universe.Id, universe.Name, universe.CreatedAt.UTC())
	return parseDbError(err)
}

func (r *universeRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (models.Universe, error) {
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

func (r *universeRepositoryImpl) List(ctx context.Context) ([]models.Universe, error) {
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

func (r *universeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.conn.Exec(ctx, deleteUniverseQuery, id)
	return err
}

func loadUniverseDetails(ctx context.Context, tx db.Transaction, dbUniverse mappers.DbUniverse) (models.Universe, error) {
	universe := dbUniverse.ToDomain()

	return universe, nil
}
