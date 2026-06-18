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
	createPlanetQuery = `
INSERT INTO
	planet (id, player, name, created_at, updated_at, version)
	VALUES($1, $2, $3, $4, $5, $6)`

	createPlanetHomeworldQuery = `INSERT INTO homeworld (player, planet) VALUES($1, $2)`

	createPlanetResourceQuery = `
INSERT INTO
	planet_resource (planet, resource, amount, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5)`

	createPlanetResourceStorageQuery = `
INSERT INTO
	planet_resource_storage (planet, resource, storage, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5)`

	createPlanetResourceProductionQuery = `
INSERT INTO
	planet_resource_production (planet, building, resource, production, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5, $6)`

	createPlanetBuildingQuery = `
INSERT INTO
	planet_building (planet, building, level, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5)`

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
	amount,
	created_at,
	updated_at
FROM
	planet_resource
WHERE
	planet = $1`

	listPlanetResourceStorageForPlanetQuery = `
SELECT
	resource,
	storage,
	created_at,
	updated_at
FROM
	planet_resource_storage
WHERE
	planet = $1`

	listPlanetResourceProductionForPlanetQuery = `
SELECT
	building,
	resource,
	production,
	created_at,
	updated_at
FROM
	planet_resource_production
WHERE
	planet = $1`

	listPlanetBuildingForPlanetQuery = `
SELECT
	building,
	level,
	created_at,
	updated_at
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

	deletePlanetBuildingsQuery           = `DELETE FROM planet_building WHERE planet = $1`
	deletePlanetResourceProductionsQuery = `DELETE FROM planet_resource_production WHERE planet = $1`
	deletePlanetResourceStoragesQuery    = `DELETE FROM planet_resource_storage WHERE planet = $1`
	deletePlanetResourcesQuery           = `DELETE FROM planet_resource WHERE planet = $1`
	deletePlanetHomeworldQuery           = `DELETE FROM homeworld WHERE planet = $1`
	deletePlanetQuery                    = `DELETE FROM planet WHERE id = $1`

	deletePlanetBuildingForPlayerQuery = `
DELETE FROM
	planet_building AS pbd
USING
	planet_building AS pb
	LEFT JOIN planet AS p ON pb.planet = p.id
WHERE
	pbd.planet = pb.planet
	AND pbd.building = pb.building
	AND p.player = $1`
	deletePlanetResourceProductionForPlayerQuery = `
DELETE FROM
	planet_resource_production AS prpd
USING
	planet_resource_production AS prp
	LEFT JOIN planet AS p ON prp.planet = p.id
WHERE
	prpd.planet = prp.planet
	AND prpd.resource = prp.resource
	AND p.player = $1`
	deletePlanetResourceStorageForPlayerQuery = `
DELETE FROM
	planet_resource_storage AS prsd
USING
	planet_resource_storage AS prs
	LEFT JOIN planet AS p ON prs.planet = p.id
WHERE
	prsd.planet = prs.planet
	AND prsd.resource = prs.resource
	AND p.player = $1`
	deletePlanetResourcesForPlayerQuery = `
DELETE FROM
	planet_resource AS prd
USING
	planet_resource AS pr
	LEFT JOIN planet AS p ON pr.planet = p.id
WHERE
	prd.planet = pr.planet
	AND prd.resource = pr.resource
	AND p.player = $1`
	deletePlanetHomeworldForPlayerQuery = `DELETE FROM homeworld WHERE player = $1`
	deletePlanetSqlForQuery             = `DELETE FROM planet WHERE player = $1`
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

	if planet.Homeworld {
		_, err = tx.Exec(ctx, createPlanetHomeworldQuery, planet.Player, planet.Id)
		if err != nil {
			return err
		}
	}

	for _, r := range planet.Resources {
		_, err = tx.Exec(
			ctx,
			createPlanetResourceQuery,
			planet.Id,
			r.Resource,
			r.Amount,
			r.CreatedAt,
			r.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	for _, s := range planet.Storages {
		_, err = tx.Exec(
			ctx,
			createPlanetResourceStorageQuery,
			planet.Id,
			s.Resource,
			s.Storage,
			s.CreatedAt,
			s.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	for _, p := range planet.Productions {
		_, err = tx.Exec(
			ctx,
			createPlanetResourceProductionQuery,
			planet.Id,
			p.Building,
			p.Resource,
			p.Production,
			p.CreatedAt,
			p.UpdatedAt,
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
			b.CreatedAt,
			b.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return err
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

	err = deleteBuildingActionDetailsForPlanet(ctx, tx, id)
	if err != nil {
		return err
	}

	return deletePlanetDetails(ctx, tx, id)
}

func (r *planetRepositoryImpl) DeleteForPlayer(ctx context.Context, player uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	err = deleteBuildingActionDetailsForPlayer(ctx, tx, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetBuildingForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceProductionForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourceStorageForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetResourcesForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetHomeworldForPlayerQuery, player)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deletePlanetSqlForQuery, player)
	return err
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

func deletePlanetDetails(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetBuildingsQuery, id)
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
