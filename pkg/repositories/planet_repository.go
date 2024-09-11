package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetRepository interface {
	Create(ctx context.Context, tx db.Transaction, planet persistence.Planet) (persistence.Planet, error)
	Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Planet, error)
	List(ctx context.Context, tx db.Transaction) ([]persistence.Planet, error)
	ListForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) ([]persistence.Planet, error)
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

const createPlanetSqlTemplate = "INSERT INTO planet (id, player, name, created_at) VALUES($1, $2, $3, $4)"
const createPlanetHomeworldSqlTemplate = "INSERT INTO homeworld (player, planet) VALUES($1, $2)"

// https://stackoverflow.com/questions/4141370/sql-insert-with-select-and-hard-coded-values
const createPlanetResourcesSqlTemplate = `
INSERT INTO
	planet_resource (planet, resource, amount, production, created_at)
SELECT
	$1,
	id,
	start_amount,
	start_production,
	$2
FROM
	resource
`

const createPlanetBuildingsSqlTemplate = `
INSERT INTO
	planet_building (planet, building, level, created_at)
SELECT
	$1,
	id,
	0,
	$2
FROM
	building
`

func (r *planetRepositoryImpl) Create(ctx context.Context, tx db.Transaction, planet persistence.Planet) (persistence.Planet, error) {
	_, err := tx.Exec(ctx, createPlanetSqlTemplate, planet.Id, planet.Player, planet.Name, planet.CreatedAt)
	if err != nil {
		return planet, err
	}

	if planet.Homeworld {
		_, err = tx.Exec(ctx, createPlanetHomeworldSqlTemplate, planet.Player, planet.Id)
		if err != nil {
			return planet, err
		}
	}

	_, err = tx.Exec(ctx, createPlanetResourcesSqlTemplate, planet.Id, planet.CreatedAt)
	if err != nil {
		return planet, err
	}

	_, err = tx.Exec(ctx, createPlanetBuildingsSqlTemplate, planet.Id, planet.CreatedAt)

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

func (r *planetRepositoryImpl) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Planet, error) {
	res := tx.Query(ctx, getPlanetSqlTemplate, id)
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

func (r *planetRepositoryImpl) List(ctx context.Context, tx db.Transaction) ([]persistence.Planet, error) {
	res := tx.Query(ctx, listPlanetSqlTemplate)
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

func (r *planetRepositoryImpl) ListForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) ([]persistence.Planet, error) {
	res := tx.Query(ctx, listPlanetForPlayerSqlTemplate, player)
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
const deletePlanetHomeworldSqlTemplate = "DELETE FROM homeworld WHERE planet = $1"

func (r *planetRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	affected, err := tx.Exec(ctx, deletePlanetHomeworldSqlTemplate, id)
	if err != nil {
		return err
	}
	if affected != 0 && affected != 1 {
		return errors.NewCode(db.MoreThanOneMatchingSqlRows)
	}

	affected, err = tx.Exec(ctx, deletePlanetSqlTemplate, id)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.NewCode(db.NoMatchingSqlRows)
	}
	return nil
}
