package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
)

const (
	createPlanetQuery = `
INSERT INTO
	planet (id, player, name, created_at, updated_at, version)
	VALUES ($1, $2, $3, $4, $5, $6)`

	createPlanetHomeworldQuery = `INSERT INTO homeworld (player, planet) VALUES ($1, $2)`

	createPlanetResourceQuery = `
INSERT INTO
	planet_resource (planet, resource, amount)
	VALUES ($1, $2, $3)`

	createPlanetResourceStorageQuery = `
INSERT INTO
	planet_resource_storage (planet, resource, storage)
	VALUES ($1, $2, $3)`

	createPlanetResourceProductionQuery = `
INSERT INTO
	planet_resource_production (planet, building, resource, production)
	VALUES ($1, $2, $3, $4)`

	createPlanetBuildingQuery = `
INSERT INTO
	planet_building (planet, building, level)
	VALUES ($1, $2, $3)`

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
	p.version,
	ba.id AS building_action
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
	LEFT JOIN building_action AS ba ON ba.planet = p.id
WHERE
	p.id = $1`

	listPlanetResourceForPlanetQuery = `
SELECT
	resource,
	amount
FROM
	planet_resource
WHERE
	planet = $1`

	listPlanetResourceStorageForPlanetQuery = `
SELECT
	resource,
	storage
FROM
	planet_resource_storage
WHERE
	planet = $1`

	listPlanetResourceProductionForPlanetQuery = `
SELECT
	building,
	resource,
	production
FROM
	planet_resource_production
WHERE
	planet = $1`

	listPlanetBuildingForPlanetQuery = `
SELECT
	building,
	level
FROM
	planet_building
WHERE
	planet = $1`

	listPlanetForPlayerQuery = `
SELECT
	p.id
FROM
	planet AS p
WHERE
	p.player = $1
ORDER BY
	p.created_at,
	p.name`

	updatePlanetQuery = `
UPDATE
	planet
SET
	version = $1,
	updated_at = $2
WHERE
	id = $3
	AND version = $4
	`

	updatePlanetResourcesQuery = `
UPDATE
	planet_resource
SET
	amount = $1
WHERE
	planet = $2
	AND resource = $3
	`
	updatePlanetStoragesQuery = `
UPDATE
	planet_resource_storage
SET
	storage = $1
WHERE
	planet = $2
	AND resource = $3
	`
	// https://wiki.postgresql.org/wiki/Is_distinct_from
	upsertPlanetProductionsQuery = `
INSERT INTO
	planet_resource_production (planet, resource, building, production)
	VALUES ($1, $2, $3, $4)
ON CONFLICT (planet, building, resource) DO UPDATE
SET
	production = excluded.production`
	updatePlanetBuildingsQuery = `
UPDATE
	planet_building
SET
	level = $1
WHERE
	planet = $2
	AND building = $3`

	deletePlanetBuildingsQuery           = `DELETE FROM planet_building WHERE planet = $1`
	deletePlanetResourceProductionsQuery = `DELETE FROM planet_resource_production WHERE planet = $1`
	deletePlanetResourceStoragesQuery    = `DELETE FROM planet_resource_storage WHERE planet = $1`
	deletePlanetResourcesQuery           = `DELETE FROM planet_resource WHERE planet = $1`
	deletePlanetHomeworldQuery           = `DELETE FROM homeworld WHERE planet = $1`
	deletePlanetQuery                    = `DELETE FROM planet WHERE id = $1`
)

type PlanetRepository struct {
	conn db.Connection
}

func NewPlanetRepository(conn db.Connection) *PlanetRepository {
	return &PlanetRepository{
		conn: conn,
	}
}

func (r *PlanetRepository) ListForPlayer(ctx context.Context, player uuid.UUID) ([]uuid.UUID, error) {
	return db.QueryAll[uuid.UUID](ctx, r.conn, listPlanetForPlayerQuery, player)
}

func (r *PlanetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	return deletePlanetAndDetails(ctx, tx, id)
}

func createPlanetWithDetails(ctx context.Context, tx db.Transaction, planet models.Planet) error {
	_, err := tx.Exec(
		ctx,
		createPlanetQuery,
		planet.Id,
		planet.Player,
		planet.Name,
		planet.CreatedAt.UTC(),
		planet.UpdatedAt.UTC(),
		planet.Version,
	)
	if err != nil {
		return parseDbError(err)
	}

	if planet.Homeworld {
		_, err := tx.Exec(ctx, createPlanetHomeworldQuery, planet.Player, planet.Id)
		if err != nil {
			return err
		}
	}

	for _, r := range planet.Resources {
		_, err := tx.Exec(
			ctx,
			createPlanetResourceQuery,
			planet.Id,
			r.Resource,
			r.Amount,
		)
		if err != nil {
			return err
		}
	}

	for _, s := range planet.Storages {
		_, err := tx.Exec(
			ctx,
			createPlanetResourceStorageQuery,
			planet.Id,
			s.Resource,
			s.Storage,
		)
		if err != nil {
			return err
		}
	}

	for _, p := range planet.Productions {
		_, err := tx.Exec(
			ctx,
			createPlanetResourceProductionQuery,
			planet.Id,
			p.Building,
			p.Resource,
			p.Production,
		)
		if err != nil {
			return err
		}
	}

	for _, b := range planet.Buildings {
		_, err := tx.Exec(
			ctx,
			createPlanetBuildingQuery,
			planet.Id,
			b.Building,
			b.Level,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadPlanetAndDetails(
	ctx context.Context,
	tx db.Transaction,
	id uuid.UUID,
) (models.Planet, error) {
	dbPlanet, err := db.QueryOneTx[mappers.DbPlanet](ctx, tx, getPlanetQuery, id)
	if err != nil {
		return models.Planet{}, parseDbError(err)
	}

	return loadPlanetDetails(ctx, tx, dbPlanet)
}

func loadPlanetDetails(ctx context.Context, tx db.Transaction, dbPlanet mappers.DbPlanet) (models.Planet, error) {
	planet := dbPlanet.ToDomain()

	var err error
	planet.Resources, err = db.QueryAllTx[models.PlanetResource](
		ctx,
		tx,
		listPlanetResourceForPlanetQuery,
		dbPlanet.Id,
	)
	if err != nil {
		return planet, err
	}

	planet.Storages, err = db.QueryAllTx[models.PlanetResourceStorage](
		ctx,
		tx,
		listPlanetResourceStorageForPlanetQuery,
		dbPlanet.Id,
	)
	if err != nil {
		return planet, err
	}

	planet.Productions, err = db.QueryAllTx[models.PlanetResourceProduction](
		ctx,
		tx,
		listPlanetResourceProductionForPlanetQuery,
		dbPlanet.Id,
	)
	if err != nil {
		return planet, err
	}

	planet.Buildings, err = db.QueryAllTx[models.PlanetBuilding](
		ctx,
		tx,
		listPlanetBuildingForPlanetQuery,
		dbPlanet.Id,
	)
	if err != nil {
		return planet, err
	}

	if dbPlanet.BuildingAction != nil {
		action, err := loadBuildingActionAndDetails(ctx, tx, *dbPlanet.BuildingAction)
		if err != nil {
			return planet, err
		}

		planet.BuildingAction = &action
	}

	return planet, nil
}

func updatePlanetDetails(
	ctx context.Context,
	tx db.Transaction,
	planet models.Planet,
	expectedVersion int,
) error {
	for _, r := range planet.Resources {
		affected, err := tx.Exec(
			ctx,
			updatePlanetResourcesQuery,
			r.Amount,
			planet.Id,
			r.Resource,
		)
		if err != nil {
			return err
		}
		if affected != 1 {
			return domainerrors.ErrNotFound
		}
	}

	for _, s := range planet.Storages {
		affected, err := tx.Exec(
			ctx,
			updatePlanetStoragesQuery,
			s.Storage,
			planet.Id,
			s.Resource,
		)
		if err != nil {
			return err
		}
		if affected != 1 {
			return domainerrors.ErrResourceNotFound
		}
	}

	err := recreateResourceProductions(ctx, tx, planet)
	if err != nil {
		return err
	}

	for _, b := range planet.Buildings {
		affected, err := tx.Exec(
			ctx,
			updatePlanetBuildingsQuery,
			b.Level,
			planet.Id,
			b.Building,
		)
		if err != nil {
			return err
		}
		if affected != 1 {
			return domainerrors.ErrBuildingNotFound
		}
	}

	err = recreateBuildingAction(ctx, tx, planet)
	if err != nil {
		return err
	}

	affected, err := tx.Exec(
		ctx,
		updatePlanetQuery,
		planet.Version,
		planet.UpdatedAt,
		planet.Id,
		expectedVersion,
	)
	if err != nil {
		return err
	}
	if affected != 1 {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return domainerrors.ErrOptimisticLocking
	}

	return nil
}

func deletePlanetAndDetails(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	err := deleteBuildingActionAndDetailsForPlanet(ctx, tx, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetBuildingsQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceProductionsQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceStoragesQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourcesQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetHomeworldQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetQuery, id)
	if err != nil {
		return err
	}

	return nil
}

// recreateResourceProductions deletes the resource production attached to a planet and recreate them
// completely. It allows to handle cases where a mutator function removed some building production as
// a building gets demolished.
func recreateResourceProductions(ctx context.Context, tx db.Transaction, planet models.Planet) error {
	_, err := tx.Exec(ctx, deletePlanetResourceProductionsQuery, planet.Id)
	if err != nil {
		return err
	}

	for _, p := range planet.Productions {
		affected, err := tx.Exec(
			ctx,
			upsertPlanetProductionsQuery,
			planet.Id,
			p.Resource,
			p.Building,
			p.Production,
		)
		if err != nil {
			return err
		}
		if affected != 1 {
			return domainerrors.ErrNotFound
		}
	}

	return nil
}

// recreateBuildingAction deletes the action first and recreate it completely: this allows to tackle
// situations where the mutator completed an existing action and recreated a new one.
func recreateBuildingAction(ctx context.Context, tx db.Transaction, planet models.Planet) error {
	err := deleteBuildingActionAndDetailsForPlanet(ctx, tx, planet.Id)
	if err != nil {
		return err
	}

	if planet.BuildingAction != nil {
		err = upsertBuildingActionWithDetails(ctx, tx, planet.Id, *planet.BuildingAction)
		if err != nil {
			return err
		}
	}

	return nil
}
