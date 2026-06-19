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

	listPlayerQuery = `
SELECT
	id,
	api_user,
	universe,
	name,
	created_at,
	version
FROM
	player
ORDER BY
	created_at,
	name`

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
	api_user = $1
ORDER BY
	created_at,
	name`

	deletePlayerQuery = `DELETE FROM player WHERE id = $1`
)

type playerRepositoryImpl struct {
	conn db.Connection
}

func NewPlayerRepository(conn db.Connection) drivenports.ForManagingPlayers {
	return &playerRepositoryImpl{
		conn: conn,
	}
}

func (r *playerRepositoryImpl) Create(
	ctx context.Context,
	player models.Player,
	homeworld models.Planet,
) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = r.conn.Exec(ctx, createPlayerQuery, player.Id, player.ApiUser, player.Universe, player.Name, player.CreatedAt)
	if err != nil {
		return parseDbError(err)
	}

	err = createPlanetWithDetails(ctx, tx, homeworld)
	if err != nil {
		return parseDbError(err)
	}

	return nil
}

func (r *playerRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (models.Player, error) {
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

func (r *playerRepositoryImpl) List(ctx context.Context) ([]models.Player, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbPlayers, err := db.QueryAllTx[mappers.DbPlayer](ctx, tx, listPlayerQuery)
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

func (r *playerRepositoryImpl) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]models.Player, error) {
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

func (r *playerRepositoryImpl) Delete(ctx context.Context, player models.Player) error {
	_, err := r.conn.Exec(ctx, deletePlayerQuery, player.Id)
	return err
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
