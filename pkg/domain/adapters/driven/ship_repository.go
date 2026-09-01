package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

const (
	getShipQuery = `
SELECT
	id,
	name,
	created_at
FROM
	ship
WHERE
	id = $1`

	listShipQuery = `
SELECT
	id,
	name,
	created_at
FROM
	ship
ORDER BY
	created_at,
	name`

	listShipCostForShipQuery = `
SELECT
	sc.resource,
	sc.cost,
	r.hours_per_unit AS build_time_hours_per_unit
FROM
	ship_cost AS sc
	INNER JOIN resource_metabolization_rate_shipyard AS r ON r.resource = sc.resource
WHERE
	sc.ship = $1`
)

type ShipRepository struct {
	conn db.Connection
}

func NewShipRepository(conn db.Connection) *ShipRepository {
	return &ShipRepository{
		conn: conn,
	}
}

func (r *ShipRepository) Get(ctx context.Context, id uuid.UUID) (models.Ship, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return models.Ship{}, err
	}
	defer tx.Close(ctx)

	dbShip, err := db.QueryOneTx[mappers.DbShip](ctx, tx, getShipQuery, id)
	if err != nil {
		return models.Ship{}, parseDbError(err)
	}

	return loadShipDetails(ctx, tx, dbShip)
}

func loadShips(ctx context.Context, tx db.Transaction) ([]models.Ship, error) {
	dbShips, err := db.QueryAllTx[mappers.DbShip](ctx, tx, listShipQuery)
	if err != nil {
		return nil, err
	}

	ships := make([]models.Ship, 0, len(dbShips))
	for id := range dbShips {
		ship, err := loadShipDetails(ctx, tx, dbShips[id])
		if err != nil {
			return nil, err
		}

		ships = append(ships, ship)
	}

	return ships, nil
}

func loadShipDetails(ctx context.Context, tx db.Transaction, dbShip mappers.DbShip) (models.Ship, error) {
	ship := dbShip.ToDomain()

	var err error
	ship.Costs, err = db.QueryAllTx[models.ShipCost](
		ctx,
		tx,
		listShipCostForShipQuery,
		dbShip.Id,
	)
	if err != nil {
		return ship, err
	}

	return ship, nil
}
