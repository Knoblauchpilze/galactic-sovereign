package domainservices

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

// AdvanceSimulationUntil executes simulation steps related to a planet. The function
// will repeatedly call the simulation step executor with the provided resolver.
// The process will stop either when an error is returned or when the resolver indicates
// that the planet was updated (returning PlanetUpdated or FleetCreated as a result).
// This process allows to progressively resolve outstanding fleet events until the
// planet is ready to be updated. All updates will be done with the provided timestamp
// meaning that the simulation will be updated up to this point in time.
func AdvanceSimulationUntil(
	ctx context.Context,
	executor drivenports.ForExecutingSimulationStep,
	planetId uuid.UUID,
	timestamp time.Time,
	resolver models.SimulationStepResolver,
) (models.SimulationStepResult, error) {
	simulate := ResolveSimulationStepThen(resolver)

	for {
		result, err := executor.Execute(ctx, planetId, timestamp, simulate)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		switch result.Resolution {
		case models.FleetEventResolved, models.FleetDeleted:
			continue
		default:
			return result, nil
		}
	}
}

// AdvanceSimulationForDeleteUntil executes simulation steps related to a planet leading
// to its deletion. The function will repeatedly call the simulation step executor with
// the provided resolver.
// The process will stop either when an error is returned or when the resolver indicates
// that the planet was deleted (returning PlanetDeleted as a result).
// This process allows to progressively resolve outstanding fleet events until the
// planet is ready to be deleted. All updates will be done with the provided timestamp
// meaning that the simulation will be updated up to this point in time.
func AdvanceSimulationForDeleteUntil(
	ctx context.Context,
	executor drivenports.ForExecutingSimulationStep,
	planetId uuid.UUID,
	timestamp time.Time,
	resolver models.SimulationStepResolver,
) (models.SimulationStepResult, error) {
	simulate := ResolveSimulationStepThen(resolver)

	for {
		result, err := executor.ExecuteForDelete(ctx, planetId, timestamp, simulate)
		if err != nil {
			return models.SimulationStepResult{}, err
		}

		switch result.Resolution {
		case models.FleetEventResolved, models.FleetDeleted:
			continue
		default:
			return result, nil
		}
	}
}
