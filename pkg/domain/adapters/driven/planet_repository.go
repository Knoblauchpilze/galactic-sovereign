package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

const (
	createPlanetQuery = `
INSERT INTO
	planet (id, player, name, created_at, updated_at, version)
	VALUES($1, $2, $3, $4, $5, $6)`

	createPlanetHomeworldQuery = `INSERT INTO homeworld (player, planet) VALUES($1, $2)`

	createPlanetResourceQuery = `
INSERT INTO
	planet_resource (planet, resource, amount)
	VALUES($1, $2, $3)`

	createPlanetResourceStorageQuery = `
INSERT INTO
	planet_resource_storage (planet, resource, storage)
	VALUES($1, $2, $3)`

	createPlanetResourceProductionQuery = `
INSERT INTO
	planet_resource_production (planet, building, resource, production)
	VALUES($1, $2, $3, $4)`

	createPlanetBuildingQuery = `
INSERT INTO
	planet_building (planet, building, level)
	VALUES($1, $2, $3)`

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
	p.version,
	ba.id AS building_action
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
	LEFT JOIN building_action AS ba ON ba.planet = p.id
ORDER BY
	p.created_at,
	p.name`

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
	p.version,
	ba.id AS building_action
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
	LEFT JOIN building_action AS ba ON ba.planet = p.id
WHERE
	p.player = $1
ORDER BY
	p.created_at,
	p.name`

	updatePlanetResourcesQuery = `
UPDATE
	planet_resource
SET
	amount = $1
WHERE
	planet = $2
	AND resource = $3
	`

	deletePlanetBuildingsQuery           = `DELETE FROM planet_building WHERE planet = $1`
	deletePlanetResourceProductionsQuery = `DELETE FROM planet_resource_production WHERE planet = $1`
	deletePlanetResourceStoragesQuery    = `DELETE FROM planet_resource_storage WHERE planet = $1`
	deletePlanetResourcesQuery           = `DELETE FROM planet_resource WHERE planet = $1`
	deletePlanetHomeworldQuery           = `DELETE FROM homeworld WHERE planet = $1`
	deletePlanetQuery                    = `DELETE FROM planet WHERE id = $1`
)

type planetRepositoryImpl struct {
	conn       db.Connection
	actionRepo *buildingActionRepositoryImpl
}

func NewPlanetRepository(conn db.Connection) drivenports.ForManagingPlanets {
	return &planetRepositoryImpl{
		conn:       conn,
		actionRepo: &buildingActionRepositoryImpl{conn: conn},
	}
}

func (r *planetRepositoryImpl) Create(ctx context.Context, planet models.Planet) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	err = createPlanetWithDetails(ctx, tx, planet)
	if err != nil {
		return parseDbError(err)
	}

	return nil
}

func (r *planetRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (models.Planet, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return models.Planet{}, err
	}
	defer tx.Close(ctx)

	dbPlanet, err := db.QueryOneTx[mappers.DbPlanet](ctx, tx, getPlanetQuery, id)
	if err != nil {
		return models.Planet{}, parseDbError(err)
	}

	return loadPlanetDetails(ctx, tx, dbPlanet)
}

func (r *planetRepositoryImpl) List(ctx context.Context) ([]models.Planet, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbPlanets, err := db.QueryAllTx[mappers.DbPlanet](ctx, tx, listPlanetQuery)
	if err != nil {
		return nil, err
	}

	planets := make([]models.Planet, 0, len(dbPlanets))
	for id := range dbPlanets {
		planet, err := loadPlanetDetails(ctx, tx, dbPlanets[id])
		if err != nil {
			return nil, err
		}

		planets = append(planets, planet)
	}

	return planets, nil
}

func (r *planetRepositoryImpl) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbPlanets, err := db.QueryAllTx[mappers.DbPlanet](ctx, tx, listPlanetForPlayerQuery, player)
	if err != nil {
		return nil, err
	}

	planets := make([]models.Planet, 0, len(dbPlanets))
	for id := range dbPlanets {
		planet, err := loadPlanetDetails(ctx, tx, dbPlanets[id])
		if err != nil {
			return nil, err
		}

		planets = append(planets, planet)
	}

	return planets, nil
}

func (r *planetRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	return deletePlanetDetails(ctx, tx, id)
}

func createPlanetWithDetails(ctx context.Context, tx db.Transaction, planet models.Planet) error {
	_, err := tx.Exec(
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

	return planet, nil
}

func updatePlanetDetails(ctx context.Context, tx db.Transaction, planet models.Planet) error {
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

	return nil
}

func deletePlanetDetails(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	err := deleteBuildingActionDetailsForPlanet(ctx, tx, id)
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
