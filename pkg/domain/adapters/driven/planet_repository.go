package driven

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

const (
	createPlanetQuery = `
INSERT INTO
	planet (id, player, name, created_at, updated_at, version)
	VALUES($1, $2, $3, $4, $5, $6)`

	createPlanetHomeworldQuery = `INSERT INTO homeworld (player, planet) VALUES($1, $2)`

	getPlanetQuery = `
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at,
	p.version
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
WHERE
	id = $1`

	listPlanetQuery = `
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at,
	p.version
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id`

	listPlanetForPlayerQuery = `
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at,
	p.version
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
WHERE
	p.player = $1`

	deletePlanetHomeworldQuery = `DELETE FROM homeworld WHERE planet = $1`

	deletePlanetQuery = `DELETE FROM planet WHERE id = $1`

	deletePlanetHomeworldForPlayerQuery = `DELETE FROM homeworld WHERE player = $1`

	deletePlanetSqlForQuery = `DELETE FROM planet WHERE player = $1`
)

type planetRepositoryImpl struct {
	conn db.Connection
}

func NewPlanetRepository(conn db.Connection) driven.ForManagingPlanets {
	return &planetRepositoryImpl{
		conn: conn,
	}
}

func (r *planetRepositoryImpl) Create(ctx context.Context, planet models.Planet) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(
		ctx,
		createPlanetQuery,
		planet.Id,
		planet.Player,
		planet.Name,
		planet.CreatedAt,
		planet.UpdatedAt,
		planet.Version,
	)
	if err != nil {
		return err
	}

	if !planet.Homeworld {
		return nil
	}

	_, err = tx.Exec(ctx, createPlanetHomeworldQuery, planet.Player, planet.Id)

	return err
}

func (r *planetRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (models.Planet, error) {
	return db.QueryOne[models.Planet](ctx, r.conn, getPlanetQuery, id)
}

func (r *planetRepositoryImpl) List(ctx context.Context) ([]models.Planet, error) {
	return db.QueryAll[models.Planet](ctx, r.conn, listPlanetQuery)
}

func (r *planetRepositoryImpl) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	return db.QueryAll[models.Planet](ctx, r.conn, listPlanetForPlayerQuery, player)
}

func (r *planetRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deletePlanetHomeworldQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetQuery, id)
	return err
}

func (r *planetRepositoryImpl) DeleteForPlayer(ctx context.Context, player uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deletePlanetHomeworldForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetSqlForQuery, player)
	return err
}
