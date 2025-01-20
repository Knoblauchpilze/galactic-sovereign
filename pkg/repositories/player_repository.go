package repositories

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlayerRepository interface {
	Create(ctx context.Context, tx db.Transaction, player persistence.Player) (persistence.Player, error)
	Get(ctx context.Context, id uuid.UUID) (persistence.Player, error)
	List(ctx context.Context) ([]persistence.Player, error)
	ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]persistence.Player, error)
	Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error
}

type playerRepositoryImpl struct {
	conn db.Connection
}

func NewPlayerRepository(conn db.Connection) PlayerRepository {
	return &playerRepositoryImpl{
		conn: conn,
	}
}

const createPlayerSqlTemplate = `
INSERT INTO
	player (id, api_user, universe, name, created_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING updated_at`

func (r *playerRepositoryImpl) Create(ctx context.Context, tx db.Transaction, player persistence.Player) (persistence.Player, error) {
	updatedAt, err := db.QueryOneTx[time.Time](ctx, tx, createPlayerSqlTemplate, player.Id, player.ApiUser, player.Universe, player.Name, player.CreatedAt)
	player.UpdatedAt = updatedAt
	return player, err
}

const getPlayerSqlTemplate = `
SELECT
	id,
	api_user,
	universe,
	name,
	created_at,
	updated_at,
	version
FROM
	player
WHERE
	id = $1`

func (r *playerRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.Player, error) {
	return db.QueryOne[persistence.Player](ctx, r.conn, getPlayerSqlTemplate, id)
}

const listPlayerSqlTemplate = `
SELECT
	id,
	api_user,
	universe,
	name,
	created_at,
	updated_at,
	version
FROM
	player`

func (r *playerRepositoryImpl) List(ctx context.Context) ([]persistence.Player, error) {
	return db.QueryAll[persistence.Player](ctx, r.conn, listPlayerSqlTemplate)
}

const listPlayerForApiUserSqlTemplate = `
SELECT
	id,
	api_user,
	universe,
	name,
	created_at,
	updated_at,
	version
FROM
	player
WHERE
	api_user = $1`

func (r *playerRepositoryImpl) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]persistence.Player, error) {
	return db.QueryAll[persistence.Player](ctx, r.conn, listPlayerForApiUserSqlTemplate, apiUser)
}

const deletePlayerSqlTemplate = `DELETE FROM player WHERE id = $1`

func (r *playerRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlayerSqlTemplate, id)
	return err
}
