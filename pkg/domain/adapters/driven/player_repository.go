package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

const (
	createPlayerQuery = `
INSERT INTO
	player (id, api_user, universe, name, created_at)
	VALUES ($1, $2, $3, $4, $5)`

	// TODO: Should be an INNER JOIN as there is always a homeworld
	getPlayerQuery = `
SELECT
	p.id,
	p.api_user,
	p.universe,
	p.name,
	p.created_at,
	p.version,
	h.planet AS homeworld
FROM
	player AS p
	LEFT JOIN homeworld AS h ON h.player = p.id
WHERE
	p.id = $1`

	listPlanetIdsForPlayerQuery = `
SELECT
	id
FROM
	planet
WHERE
	player = $1
ORDER BY
	created_at,
	name`

	// TODO: Should be an INNER JOIN as there is always a homeworld
	listPlayerForApiUserQuery = `
SELECT
	p.id,
	p.api_user,
	p.universe,
	p.name,
	p.created_at,
	p.version,
	h.planet AS homeworld
FROM
	player AS p
	LEFT JOIN homeworld AS h ON h.player = p.id
WHERE
	p.api_user = $1
ORDER BY
	p.created_at,
	p.name`

	deletePlayerQuery = `DELETE FROM player WHERE id = $1`
)

type PlayerRepository struct {
	conn db.Connection
}

func NewPlayerRepository(conn db.Connection) *PlayerRepository {
	return &PlayerRepository{
		conn: conn,
	}
}

func (r *PlayerRepository) Create(
	ctx context.Context,
	player models.Player,
	homeworld models.Planet,
) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = r.conn.Exec(
		ctx,
		createPlayerQuery,
		player.Id,
		player.ApiUser,
		player.Universe,
		player.Name,
		player.CreatedAt.UTC(),
	)
	if err != nil {
		return parseDbError(err)
	}

	err = createPlanetWithDetails(ctx, tx, homeworld)
	if err != nil {
		return parseDbError(err)
	}

	return nil
}

func (r *PlayerRepository) Get(ctx context.Context, id uuid.UUID) (models.Player, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return models.Player{}, err
	}
	defer tx.Close(ctx)

	dbPlayer, err := db.QueryOneTx[mappers.DbPlayer](ctx, tx, getPlayerQuery, id)
	if err != nil {
		return models.Player{}, parseDbError(err)
	}

	return loadPlayerDetails(ctx, tx, dbPlayer)
}

func (r *PlayerRepository) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]models.Player, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbPlayers, err := db.QueryAllTx[mappers.DbPlayer](ctx, tx, listPlayerForApiUserQuery, apiUser)
	if err != nil {
		return nil, err
	}

	players := make([]models.Player, 0, len(dbPlayers))
	for id := range dbPlayers {
		player, err := loadPlayerDetails(ctx, tx, dbPlayers[id])
		if err != nil {
			return nil, err
		}

		players = append(players, player)
	}

	return players, nil
}

func (r *PlayerRepository) Delete(ctx context.Context, player models.Player) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	for _, p := range player.Planets {
		err = deletePlanetAndDetails(ctx, tx, p)
		if err != nil {
			return parseDbError(err)
		}
	}

	_, err = tx.Exec(ctx, deletePlayerQuery, player.Id)
	if err != nil {
		return parseDbError(err)
	}

	return nil
}

func loadPlayerDetails(ctx context.Context, tx db.Transaction, dbPlayer mappers.DbPlayer) (models.Player, error) {
	player := dbPlayer.ToDomain()

	var err error
	player.Planets, err = db.QueryAllTx[uuid.UUID](
		ctx,
		tx,
		listPlanetIdsForPlayerQuery,
		dbPlayer.Id,
	)
	if err != nil {
		return player, err
	}

	return player, nil
}
