package repositories

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type UniverseRepository interface {
	Create(ctx context.Context, universe persistence.Universe) (persistence.Universe, error)
	Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Universe, error)
	List(ctx context.Context) ([]persistence.Universe, error)
	Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error
}

type universeRepositoryImpl struct {
	conn db.Connection
}

func NewUniverseRepository(conn db.Connection) UniverseRepository {
	return &universeRepositoryImpl{
		conn: conn,
	}
}

const createUniverseSqlTemplate = `
INSERT INTO
	universe (id, name, created_at)
	VALUES($1, $2, $3)
	RETURNING updated_at`

func (r *universeRepositoryImpl) Create(ctx context.Context, universe persistence.Universe) (persistence.Universe, error) {
	updatedAt, err := db.QueryOne[time.Time](ctx, r.conn, createUniverseSqlTemplate, universe.Id, universe.Name, universe.CreatedAt)
	universe.UpdatedAt = updatedAt
	return universe, err
}

const getUniverseSqlTemplate = `
SELECT
	id,
	name,
	created_at,
	updated_at,
	version
FROM
	universe
WHERE
	id = $1`

func (r *universeRepositoryImpl) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Universe, error) {
	return db.QueryOneTx[persistence.Universe](ctx, tx, getUniverseSqlTemplate, id)
}

const listUniverseSqlTemplate = `
SELECT
	id,
	name,
	created_at,
	updated_at,
	version
FROM
	universe`

func (r *universeRepositoryImpl) List(ctx context.Context) ([]persistence.Universe, error) {
	return db.QueryAll[persistence.Universe](ctx, r.conn, listUniverseSqlTemplate)
}

const deleteUniverseSqlTemplate = `DELETE FROM universe WHERE id = $1`

func (r *universeRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteUniverseSqlTemplate, id)
	return err
}
