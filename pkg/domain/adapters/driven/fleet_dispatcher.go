package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

type FleetDispatcher struct {
	conn database.Connection
}

func NewFleetCreator(conn database.Connection) *FleetDispatcher {
	return &FleetDispatcher{
		conn: conn,
	}
}

func (c *FleetDispatcher) Dispatch(
	ctx context.Context,
	req request.FleetCreationRequest,
	dispatcher drivenports.FleetCreator,
) (models.Fleet, error) {
	tx, err := c.conn.BeginTx(ctx)
	if err != nil {
		return models.Fleet{}, err
	}
	defer tx.Close(ctx)

	actual, err := db.QueryOneTx[uuid.UUID](ctx, tx, lockPlanetForUpdateQuery, req.Planet)
	if err != nil {
		return models.Fleet{}, parseDbError(err)
	}
	if actual != req.Planet {
		return models.Fleet{}, domainerrors.ErrNotFound
	}

	planet, err := loadPlanetAndDetails(ctx, tx, req.Planet)
	if err != nil {
		return models.Fleet{}, err
	}

	expectedVersion := planet.Version

	fleet, err := dispatcher(&planet, req)
	if err != nil {
		// There's no point in checking the error here: it is not logged
		// and there's already an error pending.
		// nolint:errcheck
		tx.Rollback()

		return models.Fleet{}, err
	}

	_, err = saveAndReloadPlanet(ctx, tx, planet, expectedVersion)
	if err != nil {
		return models.Fleet{}, err
	}

	return fleet, nil
}
