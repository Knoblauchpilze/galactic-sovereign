package drivenadapters

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
)

type SimulationStepExecutor struct {
	conn database.Connection
}

func NewSimulationStepExecutor(conn database.Connection) *SimulationStepExecutor {
	return &SimulationStepExecutor{
		conn: conn,
	}
}

// TODO: Add tests for this
func (s *SimulationStepExecutor) Execute(
	ctx context.Context,
	planetId uuid.UUID,
	timestamp time.Time,
	simulate func(models.SimulationStep) (models.SimulationStepResult, error),
) (models.SimulationStepResult, error) {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return models.SimulationStepResult{}, err
	}
	defer tx.Close(ctx)

	actual, err := db.QueryOneTx[uuid.UUID](ctx, tx, lockPlanetForUpdateQuery, planetId)
	if err != nil {
		return models.SimulationStepResult{}, parseDbError(err)
	}
	if actual != planetId {
		return models.SimulationStepResult{}, domainerrors.ErrNotFound
	}

	planet, err := loadPlanetAndDetails(ctx, tx, planetId)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	expectedVersion := planet.Version

	step := models.SimulationStep{
		Timestamp:      timestamp,
		Planet:         planet,
		NextFleetEvent: nil,
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
		out, err := saveAndReloadPlanet(ctx, tx, *result.Planet, expectedVersion)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		result.Planet = &out
	}

	if result.Fleet != nil {
		// TODO: Add tests for this
		return models.SimulationStepResult{}, errors.ErrNotImplemented
	}

	return result, nil
}
