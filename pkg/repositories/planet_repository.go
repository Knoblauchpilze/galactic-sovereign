package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetRepository interface {
	Create(ctx context.Context, tx db.Transaction, planet persistence.Planet) (persistence.Planet, error)
	Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Planet, error)
	List(ctx context.Context, tx db.Transaction) ([]persistence.Planet, error)
	ListForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) ([]persistence.Planet, error)
	Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error
	DeleteForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) error
}

type planetRepositoryImpl struct {
	conn db.Connection
}

func NewPlanetRepository(conn db.Connection) PlanetRepository {
	return &planetRepositoryImpl{
		conn: conn,
	}
}

const createPlanetSqlTemplate = `
INSERT INTO
	planet (id, player, name, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5)`
const createPlanetHomeworldSqlTemplate = `INSERT INTO homeworld (player, planet) VALUES($1, $2)`

// https://stackoverflow.com/questions/4141370/sql-insert-with-select-and-hard-coded-values
const createPlanetResourcesSqlTemplate = `
INSERT INTO
	planet_resource (planet, resource, amount, created_at, updated_at)
SELECT
	$1,
	id,
	start_amount,
	$2,
	$3
FROM
	resource`

const createPlanetResourceProductionsSqlTemplate = `
INSERT INTO
	planet_resource_production (planet, resource, production, created_at, updated_at)
SELECT
	$1,
	id,
	start_production,
	$2,
	$3
FROM
	resource`

const createPlanetResourceStoragesSqlTemplate = `
INSERT INTO
	planet_resource_storage (planet, resource, storage, created_at, updated_at)
SELECT
	$1,
	id,
	start_storage,
	$2,
	$3
FROM
	resource`

const createPlanetBuildingsSqlTemplate = `
INSERT INTO
	planet_building (planet, building, level, created_at, updated_at)
SELECT
	$1,
	id,
	0,
	$2,
	$3
FROM
	building`

func (r *planetRepositoryImpl) Create(ctx context.Context, tx db.Transaction, planet persistence.Planet) (persistence.Planet, error) {
	_, err := tx.Exec(ctx, createPlanetSqlTemplate, planet.Id, planet.Player, planet.Name, planet.CreatedAt, planet.CreatedAt)
	planet.UpdatedAt = planet.CreatedAt
	if err != nil {
		return planet, err
	}

	if planet.Homeworld {
		_, err = tx.Exec(ctx, createPlanetHomeworldSqlTemplate, planet.Player, planet.Id)
		if err != nil {
			return planet, err
		}
	}

	_, err = tx.Exec(ctx, createPlanetResourcesSqlTemplate, planet.Id, planet.CreatedAt, planet.UpdatedAt)
	if err != nil {
		return planet, err
	}

	_, err = tx.Exec(ctx, createPlanetResourceProductionsSqlTemplate, planet.Id, planet.CreatedAt, planet.UpdatedAt)
	if err != nil {
		return planet, err
	}

	_, err = tx.Exec(ctx, createPlanetResourceStoragesSqlTemplate, planet.Id, planet.CreatedAt, planet.UpdatedAt)
	if err != nil {
		return planet, err
	}

	_, err = tx.Exec(ctx, createPlanetBuildingsSqlTemplate, planet.Id, planet.CreatedAt, planet.UpdatedAt)

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
	id = $1`

func (r *planetRepositoryImpl) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Planet, error) {
	return db.QueryOneTx[persistence.Planet](ctx, tx, getPlanetSqlTemplate, id)
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
	LEFT JOIN homeworld AS h ON h.planet = p.id`

func (r *planetRepositoryImpl) List(ctx context.Context, tx db.Transaction) ([]persistence.Planet, error) {
	return db.QueryAllTx[persistence.Planet](ctx, tx, listPlanetSqlTemplate)
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
	p.player = $1`

func (r *planetRepositoryImpl) ListForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) ([]persistence.Planet, error) {
	return db.QueryAllTx[persistence.Planet](ctx, tx, listPlanetForPlayerSqlTemplate, player)
}

const deletePlanetBuildingsSqlTemplate = `DELETE FROM planet_building WHERE planet = $1`
const deletePlanetResourceStoragesSqlTemplate = `DELETE FROM planet_resource_storage WHERE planet = $1`
const deletePlanetResourceProductionsSqlTemplate = `DELETE FROM planet_resource_production WHERE planet = $1`
const deletePlanetResourcesSqlTemplate = `DELETE FROM planet_resource WHERE planet = $1`
const deletePlanetHomeworldSqlTemplate = `DELETE FROM homeworld WHERE planet = $1`
const deletePlanetSqlTemplate = `DELETE FROM planet WHERE id = $1`

func (r *planetRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetBuildingsSqlTemplate, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceStoragesSqlTemplate, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceProductionsSqlTemplate, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourcesSqlTemplate, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetHomeworldSqlTemplate, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetSqlTemplate, id)
	return err
}

const deletePlanetBuildingForPlayerSqlTemplate = `
DELETE FROM
	planet_building AS pbd
USING
	planet_building AS pb
	LEFT JOIN planet AS p ON pb.planet = p.id
WHERE
	pbd.planet = pb.planet
	AND pbd.building = pb.building
	AND p.player = $1`
const deletePlanetResourceStorageForPlayerSqlTemplate = `
DELETE FROM
	planet_resource_storage AS prsd
USING
	planet_resource_storage AS prs
	LEFT JOIN planet AS p ON prs.planet = p.id
WHERE
	prsd.planet = prs.planet
	AND prsd.resource = prs.resource
	AND p.player = $1`
const deletePlanetResourceProductionForPlayerSqlTemplate = `
DELETE FROM
	planet_resource_production AS prpd
USING
	planet_resource_production AS prp
	LEFT JOIN planet AS p ON prp.planet = p.id
WHERE
	prpd.planet = prp.planet
	AND prpd.resource = prp.resource
	AND p.player = $1`
const deletePlanetResourceForPlayerSqlTemplate = `
DELETE FROM
	planet_resource AS prd
USING
	planet_resource AS pr
	LEFT JOIN planet AS p ON pr.planet = p.id
WHERE
	prd.planet = pr.planet
	AND prd.resource = pr.resource
	AND p.player = $1`
const deletePlanetHomeworldForPlayerSqlTemplate = `DELETE FROM homeworld WHERE player = $1`
const deletePlanetSqlForPlayerTemplate = `DELETE FROM planet WHERE player = $1`

func (r *planetRepositoryImpl) DeleteForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetBuildingForPlayerSqlTemplate, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceStorageForPlayerSqlTemplate, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceProductionForPlayerSqlTemplate, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceForPlayerSqlTemplate, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetHomeworldForPlayerSqlTemplate, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetSqlForPlayerTemplate, player)
	return err
}
