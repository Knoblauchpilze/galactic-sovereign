package drivenports

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

// ForExecutingSimulationStep loads a planet's simulation state at a timestamp,
// passes it to a simulation function, and returns the resulting state changes.
// The simulation state includes the planet itself but also the earliest unresolved
// fleet event across the universe the planet belongs to. Both are provided to
// the decide function which should resolve the simulation step.
// The result includes the new states of the planet and fleet. Both values will
// be persisted back to the database.
// TODO: Probably the return of the function should be:
//   - event resolved/fleet resolved boolean
//   - the planet when it was updated
//
// Could also be the SimulationStep returning it -> keeps the logic
// in the domain
type ForExecutingSimulationStep interface {
	Execute(
		ctx context.Context,
		planetId uuid.UUID,
		timestamp time.Time,
		simulate func(models.SimulationStep) (models.SimulationStepResult, error),
	) (models.SimulationStepResult, error)
}
