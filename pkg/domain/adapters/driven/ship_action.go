package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

const (
	upsertShipActionQuery = `
INSERT INTO
	ship_action (id, planet, ship, count, created_at, next_completion_at, unit_completion_time)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE
SET
	count = excluded.count,
	next_completion_at = excluded.next_completion_at,
	unit_completion_time = excluded.unit_completion_time`

	listShipActionForPlanetQuery = `
SELECT
	id,
	ship,
	count,
	created_at,
	next_completion_at,
	unit_completion_time
FROM
	ship_action
WHERE
	planet = $1`

	deleteShipActionForPlanetQuery = `DELETE FROM ship_action WHERE planet = $1`
)

func upsertShipActionWithDetails(
	ctx context.Context,
	tx db.Transaction,
	planet uuid.UUID,
	action models.ShipAction,
) error {
	_, err := tx.Exec(
		ctx,
		upsertShipActionQuery,
		action.Id,
		planet,
		action.Ship,
		action.Count,
		action.CreatedAt,
		action.NextCompletionAt,
		action.UnitCompletionTime,
	)
	if err != nil {
		return err
	}

	return nil
}

func loadShipActionAndDetailsForPlanet(
	ctx context.Context,
	tx db.Transaction,
	planet uuid.UUID,
) ([]models.ShipAction, error) {
	dbShipActions, err := db.QueryAllTx[mappers.DbShipAction](
		ctx,
		tx,
		listShipActionForPlanetQuery,
		planet,
	)
	if err != nil {
		return []models.ShipAction{}, err
	}

	out := []models.ShipAction{}

	for _, dbShipAction := range dbShipActions {
		out = append(out, dbShipAction.ToDomain())
	}

	return out, nil
}

func deleteShipActionAndDetailsForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteShipActionForPlanetQuery, planet)
	if err != nil {
		return err
	}

	return nil
}
