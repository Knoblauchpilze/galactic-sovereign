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
	createUniverseQuery = `
INSERT INTO
	universe (id, name, created_at)
	VALUES ($1, $2, $3)`

	createUniverseTopologyQuery = `
INSERT INTO
	universe_topology (universe, galaxies, solar_systems, orbits)
	VALUES ($1, $2, $3, $4)`

	getUniverseQuery = `
SELECT
	u.id,
	u.name,
	u.created_at,
	u.version,
	ut.galaxies,
	ut.solar_systems,
	ut.orbits
FROM
	universe AS u
	INNER JOIN universe_topology AS ut ON ut.universe = u.id
WHERE
	u.id = $1`

	listResourceQuery = `
SELECT
	id,
	name,
	start_amount,
	start_production,
	start_storage,
	build_time_hours_per_unit,
	created_at
FROM
	resource
ORDER BY
	created_at,
	resource`

	listUsedCoordinateQuery = `
SELECT
	galaxy,
	solar_system,
	position
FROM
	planet_coordinate
WHERE
	universe = $1`

	listUniverseQuery = `
SELECT
	u.id,
	u.name,
	u.created_at,
	u.version,
	ut.galaxies,
	ut.solar_systems,
	ut.orbits
FROM
	universe AS u
	INNER JOIN universe_topology AS ut ON ut.universe = u.id
ORDER BY
	u.created_at,
	u.name`

	deleteUniverseTopologyQuery = `DELETE FROM universe_topology WHERE universe = $1`
	deleteUniverseQuery         = `DELETE FROM universe WHERE id = $1`
)

type UniverseRepository struct {
	conn db.Connection
}

func NewUniverseRepository(conn db.Connection) *UniverseRepository {
	return &UniverseRepository{
		conn: conn,
	}
}

func (r *UniverseRepository) Create(ctx context.Context, universe models.Universe) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = r.conn.Exec(ctx, createUniverseQuery, universe.Id, universe.Name, universe.CreatedAt.UTC())
	if err != nil {
		return parseDbError(err)
	}

	_, err = r.conn.Exec(
		ctx,
		createUniverseTopologyQuery,
		universe.Id,
		universe.Topology.Galaxies,
		universe.Topology.SolarSystems,
		universe.Topology.Orbits,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UniverseRepository) Get(ctx context.Context, id uuid.UUID) (models.Universe, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return models.Universe{}, err
	}
	defer tx.Close(ctx)

	dbUniverse, err := db.QueryOneTx[mappers.DbUniverse](ctx, tx, getUniverseQuery, id)
	if err != nil {
		return models.Universe{}, parseDbError(err)
	}

	return loadUniverseDetails(ctx, tx, dbUniverse)
}

func (r *UniverseRepository) List(ctx context.Context) ([]models.Universe, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close(ctx)

	dbUniverses, err := db.QueryAllTx[mappers.DbUniverse](ctx, tx, listUniverseQuery)
	if err != nil {
		return nil, err
	}

	universes := make([]models.Universe, 0, len(dbUniverses))
	for id := range dbUniverses {
		universe, err := loadUniverseDetails(ctx, tx, dbUniverses[id])
		if err != nil {
			return nil, err
		}

		universes = append(universes, universe)
	}

	return universes, nil
}

func (r *UniverseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deleteUniverseTopologyQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteUniverseQuery, id)
	if err != nil {
		outErr := parseDbError(err)
		// This error is returned when a player is still registered in a universe.
		// As removing a universe is most likely a destructive operation, no cascade
		// is implemented and so the error is returned as 'universe not empty'.
		if outErr == domainerrors.ErrUniverseNotFound {
			return domainerrors.ErrUniverseIsNotEmpty
		}

		return outErr
	}

	return nil
}

func loadUniverseDetails(ctx context.Context, tx db.Transaction, dbUniverse mappers.DbUniverse) (models.Universe, error) {
	universe := dbUniverse.ToDomain()

	var err error
	universe.Resources, err = db.QueryAllTx[models.Resource](
		ctx,
		tx,
		listResourceQuery,
	)
	if err != nil {
		return universe, err
	}

	universe.Buildings, err = loadBuildings(ctx, tx)
	if err != nil {
		return universe, err
	}

	universe.OccupancyMap, err = loadOccupancyMap(ctx, tx, universe.Id, universe.Topology)
	if err != nil {
		return universe, err
	}

	return universe, nil
}

func loadOccupancyMap(
	ctx context.Context,
	tx db.Transaction,
	universe uuid.UUID,
	topology models.UniverseTopology,
) (models.OccupancyMap, error) {
	occupancy := models.OccupancyMap{
		Topology:  topology,
		UsedSlots: make(map[models.Coordinate]struct{}),
	}

	slots, err := db.QueryAllTx[models.Coordinate](ctx, tx, listUsedCoordinateQuery, universe)
	if err != nil {
		return models.OccupancyMap{}, nil
	}

	for _, c := range slots {
		occupancy.UsedSlots[c] = struct{}{}
	}

	return occupancy, nil
}
