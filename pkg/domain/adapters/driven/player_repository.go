package driven

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

const (
	createPlayerQuery = `
INSERT INTO
	player (id, api_user, universe, name, created_at)
	VALUES ($1, $2, $3, $4, $5)`

	getPlayerQuery = `
SELECT
	id,
	api_user,
	universe,
	name,
	created_at,
	version
FROM
	player
WHERE
	id = $1`

	listPlayerQuery = `
SELECT
	id,
	api_user,
	universe,
	name,
	created_at,
	version
FROM
	player`

	listPlayerForApiUserQuery = `
SELECT
	id,
	api_user,
	universe,
	name,
	created_at,
	version
FROM
	player
WHERE
	api_user = $1`

	deletePlayerQuery = `DELETE FROM player WHERE id = $1`
)

type playerRepositoryImpl struct {
	conn db.Connection
}

func NewPlayerRepository(conn db.Connection) driven.ForManagingPlayers {
	return &playerRepositoryImpl{
		conn: conn,
	}
}

func (r *playerRepositoryImpl) Create(ctx context.Context, player models.Player) error {
	_, err := r.conn.Exec(ctx, createPlayerQuery, player.Id, player.ApiUser, player.Universe, player.Name, player.CreatedAt)
	return err
}

func (r *playerRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (models.Player, error) {
	return db.QueryOne[models.Player](ctx, r.conn, getPlayerQuery, id)
}

func (r *playerRepositoryImpl) List(ctx context.Context) ([]models.Player, error) {
	return db.QueryAll[models.Player](ctx, r.conn, listPlayerQuery)
}

func (r *playerRepositoryImpl) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]models.Player, error) {
	return db.QueryAll[models.Player](ctx, r.conn, listPlayerForApiUserQuery, apiUser)
}

func (r *playerRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.conn.Exec(ctx, deletePlayerQuery, id)
	return err
}
