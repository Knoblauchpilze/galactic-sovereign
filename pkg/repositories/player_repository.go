package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
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
	conn db.ConnectionPool
}

func NewPlayerRepository(conn db.ConnectionPool) PlayerRepository {
	return &playerRepositoryImpl{
		conn: conn,
	}
}

const createPlayerSqlTemplate = "INSERT INTO player (id, api_user, universe, name, created_at) VALUES($1, $2, $3, $4, $5)"

func (r *playerRepositoryImpl) Create(ctx context.Context, tx db.Transaction, player persistence.Player) (persistence.Player, error) {
	_, err := tx.Exec(ctx, createPlayerSqlTemplate, player.Id, player.ApiUser, player.Universe, player.Name, player.CreatedAt)
	if err != nil && duplicatedKeySqlErrorRegexp.MatchString(err.Error()) {
		return persistence.Player{}, errors.NewCode(db.DuplicatedKeySqlKey)
	}

	return player, err
}

const getPlayerSqlTemplate = "SELECT id, api_user, universe, name, created_at, updated_at, version FROM player WHERE id = $1"

func (r *playerRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.Player, error) {
	res := r.conn.Query(ctx, getPlayerSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.Player{}, err
	}

	var out persistence.Player
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.ApiUser, &out.Universe, &out.Name, &out.CreatedAt, &out.UpdatedAt, &out.Version)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.Player{}, err
	}

	return out, nil
}

const listPlayerSqlTemplate = "SELECT id, api_user, universe, name, created_at, updated_at, version FROM player"

func (r *playerRepositoryImpl) List(ctx context.Context) ([]persistence.Player, error) {
	res := r.conn.Query(ctx, listPlayerSqlTemplate)
	if err := res.Err(); err != nil {
		return []persistence.Player{}, err
	}

	var out []persistence.Player
	parser := func(rows db.Scannable) error {
		var player persistence.Player
		err := rows.Scan(&player.Id, &player.ApiUser, &player.Universe, &player.Name, &player.CreatedAt, &player.UpdatedAt, &player.Version)
		if err != nil {
			return err
		}

		out = append(out, player)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.Player{}, err
	}

	return out, nil
}

const listPlayerForApiUserSqlTemplate = "SELECT id, api_user, universe, name, created_at, updated_at, version FROM player where api_user = $1"

func (r *playerRepositoryImpl) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]persistence.Player, error) {
	res := r.conn.Query(ctx, listPlayerForApiUserSqlTemplate, apiUser)
	if err := res.Err(); err != nil {
		return []persistence.Player{}, err
	}

	var out []persistence.Player
	parser := func(rows db.Scannable) error {
		var player persistence.Player
		err := rows.Scan(&player.Id, &player.ApiUser, &player.Universe, &player.Name, &player.CreatedAt, &player.UpdatedAt, &player.Version)
		if err != nil {
			return err
		}

		out = append(out, player)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.Player{}, err
	}

	return out, nil
}

const deletePlayerSqlTemplate = "DELETE FROM player WHERE id = $1"

func (r *playerRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	affected, err := tx.Exec(ctx, deletePlayerSqlTemplate, id)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.NewCode(db.NoMatchingSqlRows)
	}
	return nil
}
