package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

const (
	createFleetQuery = `
INSERT INTO
	fleet (id, player, source, mission, created_at, updated_at, version)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	createFleetDestinationQuery = `
INSERT INTO
	fleet_destination (fleet, galaxy, solar_system, position, planet)
	VALUES ($1, $2, $3, $4, $5)`

	createFleetFlightQuery = `
INSERT INTO
	fleet_flight (fleet, arrival_at, return_at)
	VALUES ($1, $2, $3)`

	createFleetShipQuery = `
INSERT INTO
	fleet_ship (fleet, ship, count)
	VALUES ($1, $2, $3)`

	getFleetQuery = `
SELECT
	f.id,
	f.player,
	f.source,
	f.mission,
	fd.galaxy,
	fd.solar_system,
	fd.position,
	fd.planet AS target,
	f.created_at,
	ff.arrival_at,
	ff.return_at,
	f.updated_at,
	f.version
FROM
	fleet AS f
	LEFT JOIN fleet_destination AS fd ON fd.fleet = f.id
	LEFT JOIN fleet_flight AS ff ON ff.fleet = f.id
WHERE
	f.id = $1`

	listShipForFleetQuery = `
SELECT
	ship,
	count
FROM
	fleet_ship
WHERE
	fleet = $1`
)

func createFleetWithDetails(
	ctx context.Context,
	tx database.Transaction,
	fleet models.Fleet,
) error {
	_, err := tx.Exec(
		ctx,
		createFleetQuery,
		fleet.Id,
		fleet.Player,
		fleet.Source,
		fleet.Mission,
		fleet.CreatedAt,
		fleet.UpdatedAt,
		fleet.Version,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		createFleetDestinationQuery,
		fleet.Id,
		fleet.Destination.Coordinate.Galaxy,
		fleet.Destination.Coordinate.SolarSystem,
		fleet.Destination.Coordinate.Position,
		fleet.Destination.Target,
	)
	if err != nil {
		return parseDbError(err)
	}

	_, err = tx.Exec(
		ctx,
		createFleetFlightQuery,
		fleet.Id,
		fleet.ArrivalAt,
		fleet.ReturnAt,
	)
	if err != nil {
		return parseDbError(err)
	}

	for _, s := range fleet.Ships {
		_, err := tx.Exec(
			ctx,
			createFleetShipQuery,
			fleet.Id,
			s.Ship,
			s.Count,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadFleetAndDetails(
	ctx context.Context,
	tx database.Transaction,
	id uuid.UUID,
) (models.Fleet, error) {
	dbFleet, err := db.QueryOneTx[mappers.DbFleet](ctx, tx, getFleetQuery, id)
	if err != nil {
		return models.Fleet{}, parseDbError(err)
	}

	return loadFleetDetails(ctx, tx, dbFleet)
}

func loadFleetDetails(ctx context.Context, tx database.Transaction, dbFleet mappers.DbFleet) (models.Fleet, error) {
	fleet := dbFleet.ToDomain()

	var err error
	fleet.Ships, err = db.QueryAllTx[models.FleetShip](
		ctx,
		tx,
		listShipForFleetQuery,
		dbFleet.Id,
	)
	if err != nil {
		return fleet, err
	}

	return fleet, nil
}
