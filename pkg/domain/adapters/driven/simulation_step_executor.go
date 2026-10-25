package drivenadapters

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
)

const (
	lockFirstUnresolvedFleetEventForUpdateQuery = `
WITH target_universe AS (
	SELECT
		pl.universe
	FROM
		planet AS p
		INNER JOIN player AS pl ON pl.id = p.player
	WHERE
		p.id = $2
)
SELECT
	fe.fleet,
	fe.next_event_at,
	fe.event_type
FROM
	fleet_event AS fe
	INNER JOIN fleet AS f ON f.id = fe.fleet
	INNER JOIN planet AS p ON p.id = f.source
	INNER JOIN player AS pl ON pl.id = p.player
WHERE
	fe.next_event_at <= $1
	AND pl.universe = (SELECT universe FROM target_universe)
ORDER BY
	fe.next_event_at ASC,
	fe.fleet ASC
LIMIT 1
FOR UPDATE OF fe`

	// https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-ROWS
	// This variant allows to lock a planet and allow foreign key checks
	// typically to insert a fleet while checking that the corresponding
	// destination planet (which might be subject to another lock) still
	// exists.
	lockPlanetForNoKeyUpdateQuery = `
SELECT
	id
FROM
	planet
WHERE
	id = $1
FOR NO KEY UPDATE
`
)

type SimulationStepExecutor struct {
	conn database.Connection
}

type fleetEvent struct {
	Fleet       uuid.UUID
	NextEventAt time.Time
	EventType   mappers.DbFleetEventType
}

func NewSimulationStepExecutor(conn database.Connection) *SimulationStepExecutor {
	return &SimulationStepExecutor{
		conn: conn,
	}
}

func (s *SimulationStepExecutor) Execute(
	ctx context.Context,
	planetId uuid.UUID,
	timestamp time.Time,
	simulate models.SimulationStepResolver,
) (models.SimulationStepResult, error) {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return models.SimulationStepResult{}, err
	}
	defer tx.Close(ctx)

	event, err := resolveNextFleetEvent(ctx, tx, timestamp, planetId)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	var planet *models.Planet
	if event == nil {
		planet, err = loadAndLockForNoKeyUpdatePlanet(ctx, tx, planetId)
		if err != nil {
			return models.SimulationStepResult{}, err
		}
	}

	step := models.SimulationStep{
		Timestamp:      timestamp,
		Planet:         planet,
		NextFleetEvent: event,
	}

	var expectedPlanetVersion int
	if planet != nil {
		expectedPlanetVersion = planet.Version
	}
	var expectedFleetVersion int
	if step.NextFleetEvent != nil {
		if step.NextFleetEvent.Planet != nil {
			expectedPlanetVersion = step.NextFleetEvent.Planet.Version
		}
		expectedFleetVersion = step.NextFleetEvent.Fleet.Version
	}

	result, err := simulate(step)
	if err != nil {
		// There's no point in checking the error here: it is not logged
		// and there's already an error pending.
		// nolint:errcheck
		tx.Rollback()

		return models.SimulationStepResult{}, err
	}

	if result.Planet != nil {
		out, err := saveAndReloadPlanet(ctx, tx, *result.Planet, expectedPlanetVersion)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		result.Planet = &out
	}

	fleet, err := createUpdateOrDeleteFleet(ctx, tx, step, result, expectedFleetVersion)
	if err != nil {
		return models.SimulationStepResult{}, err
	}
	result.Fleet = fleet

	return result, nil
}

func resolveNextFleetEvent(
	ctx context.Context,
	tx database.Transaction,
	timestamp time.Time,
	planetId uuid.UUID,
) (*models.FleetEventContext, error) {
	maybeFleetEvent, err := db.QueryOneTx[fleetEvent](
		ctx,
		tx,
		lockFirstUnresolvedFleetEventForUpdateQuery,
		timestamp,
		planetId,
	)
	if err == db.ErrNoMatchingRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	fleet, err := loadFleetAndDetails(ctx, tx, maybeFleetEvent.Fleet)
	if err != nil {
		return nil, err
	}

	var planet *models.Planet
	switch maybeFleetEvent.EventType.ToDomain() {
	case models.ArrivalEvent:
		if fleet.Destination.Target != nil {
			planet, err = loadAndLockForNoKeyUpdatePlanet(ctx, tx, *fleet.Destination.Target)
		}
	case models.ReturnEvent:
		planet, err = loadAndLockForNoKeyUpdatePlanet(ctx, tx, fleet.Source)
	default:
		return nil, domainerrors.ErrUnsupportedFleetStatus
	}
	if err != nil {
		return nil, err
	}

	out := &models.FleetEventContext{
		Timestamp: maybeFleetEvent.NextEventAt,
		Fleet:     fleet,
		Planet:    planet,
	}

	return out, nil
}

func loadAndLockForNoKeyUpdatePlanet(ctx context.Context, tx database.Transaction, planetId uuid.UUID) (*models.Planet, error) {
	actual, err := db.QueryOneTx[uuid.UUID](ctx, tx, lockPlanetForNoKeyUpdateQuery, planetId)
	if err != nil {
		return nil, parseDbError(err)
	}
	if actual != planetId {
		return nil, domainerrors.ErrNotFound
	}

	planet, err := loadPlanetAndDetails(ctx, tx, planetId)
	if err != nil {
		return nil, err
	}

	return &planet, nil
}

func createUpdateOrDeleteFleet(
	ctx context.Context,
	tx database.Transaction,
	step models.SimulationStep,
	result models.SimulationStepResult,
	expectedFleetVersion int,
) (*models.Fleet, error) {
	if result.Resolution == models.FleetDeleted {
		err := deleteFleetAndDetails(ctx, tx, step.NextFleetEvent.Fleet.Id)
		if err != nil {
			return nil, err
		}
	}

	if result.Fleet == nil {
		return nil, nil
	}

	if result.Resolution == models.FleetCreated {
		err := createFleetWithDetails(ctx, tx, *result.Fleet)
		if err != nil {
			if isFleetMissingTargetPlanetErr(err) {
				return nil, domainerrors.ErrFleetDestinationInvalid
			}
			if isFleetSourcePlanetErr(err) {
				return nil, domainerrors.ErrPlanetNotFound
			}
			return nil, err
		}

		return result.Fleet, nil
	}

	out, err := saveAndReloadFleet(ctx, tx, *result.Fleet, expectedFleetVersion)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
