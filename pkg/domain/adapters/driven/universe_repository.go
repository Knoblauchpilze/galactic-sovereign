package drivenadapter

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenport "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
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
	universe`

	deleteUniverseQuery = `DELETE FROM universe WHERE id = $1`
)

type universeRepositoryImpl struct {
	conn db.Connection
}

func NewUniverseRepository(conn db.Connection) drivenport.ForManagingUniverses {
	return &universeRepositoryImpl{
		conn: conn,
	}
}

func (r *universeRepositoryImpl) Create(ctx context.Context, universe models.Universe) error {
	_, err := r.conn.Exec(ctx, createUniverseQuery, universe.Id, universe.Name, universe.CreatedAt.UTC())
	return err
}

func (r *universeRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (models.Universe, error) {
	return db.QueryOne[models.Universe](ctx, r.conn, getUniverseQuery, id)
}

func (r *universeRepositoryImpl) List(ctx context.Context) ([]models.Universe, error) {
	return db.QueryAll[models.Universe](ctx, r.conn, listUniverseQuery)
}

func (r *universeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.conn.Exec(ctx, deleteUniverseQuery, id)
	return err
}
