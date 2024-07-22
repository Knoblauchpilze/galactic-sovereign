package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetRepository interface {
	Create(ctx context.Context, planet persistence.Planet) (persistence.Planet, error)
	Get(ctx context.Context, id uuid.UUID) (persistence.Planet, error)
	List(ctx context.Context) ([]persistence.Planet, error)
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]persistence.Planet, error)
	Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error
}

type planetRepositoryImpl struct {
	conn db.ConnectionPool
}

func NewPlanetRepository(conn db.ConnectionPool) PlanetRepository {
	return &planetRepositoryImpl{
		conn: conn,
	}
}

const createPlanetSqlTemplate = "INSERT INTO planet (id, player, name, homeworld, created_at) VALUES($1, $2, $3, $4, $5)"

func (r *planetRepositoryImpl) Create(ctx context.Context, planet persistence.Planet) (persistence.Planet, error) {
	_, err := r.conn.Exec(ctx, createPlanetSqlTemplate, planet.Id, planet.Player, planet.Name, planet.Homeworld, planet.CreatedAt)
	return planet, err
}

const getPlanetSqlTemplate = `
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
WHERE
	id = $1
`

func (r *planetRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.Planet, error) {
	res := r.conn.Query(ctx, getPlanetSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.Planet{}, err
	}

	var out persistence.Planet
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Player, &out.Name, &out.Homeworld, &out.CreatedAt, &out.UpdatedAt)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.Planet{}, err
	}

	return out, nil
}

const listPlanetSqlTemplate = `
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
`

func (r *planetRepositoryImpl) List(ctx context.Context) ([]persistence.Planet, error) {
	res := r.conn.Query(ctx, listPlanetSqlTemplate)
	if err := res.Err(); err != nil {
		return []persistence.Planet{}, err
	}

	var out []persistence.Planet
	parser := func(rows db.Scannable) error {
		var planet persistence.Planet
		err := rows.Scan(&planet.Id, &planet.Player, &planet.Name, &planet.Homeworld, &planet.CreatedAt, &planet.UpdatedAt)
		if err != nil {
			return err
		}

		out = append(out, planet)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.Planet{}, err
	}

	return out, nil
}

const listPlanetForPlayerSqlTemplate = `
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
WHERE
	p.player = $1
`

func (r *planetRepositoryImpl) ListForPlayer(ctx context.Context, player uuid.UUID) ([]persistence.Planet, error) {
	res := r.conn.Query(ctx, listPlanetForPlayerSqlTemplate, player)
	if err := res.Err(); err != nil {
		return []persistence.Planet{}, err
	}

	var out []persistence.Planet
	parser := func(rows db.Scannable) error {
		var planet persistence.Planet
		err := rows.Scan(&planet.Id, &planet.Player, &planet.Name, &planet.Homeworld, &planet.CreatedAt, &planet.UpdatedAt)
		if err != nil {
			return err
		}

		out = append(out, planet)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.Planet{}, err
	}

	return out, nil
}

const deletePlanetSqlTemplate = "DELETE FROM planet WHERE id = $1"

func (r *planetRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	affected, err := tx.Exec(ctx, deletePlanetSqlTemplate, id)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.NewCode(db.NoMatchingSqlRows)
	}
	return nil
}
