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
type ForExecutingSimulationStep interface {
	// Execute runs the simulate step with the implicit expectation that the planet
	// to which it applies will not be deleted but only updated.
	Execute(
		ctx context.Context,
		planetId uuid.UUID,
		timestamp time.Time,
		simulate models.SimulationStepResolver,
	) (models.SimulationStepResult, error)

	// Execute runs the simulate step with the implicit expectation that the planet
	// to which it applies will be deleted at the end of the step. This means that
	// the simulation step result will never contain a planet (or a fleet).
	ExecuteForDelete(
		ctx context.Context,
		planetId uuid.UUID,
		timestamp time.Time,
		simulate models.SimulationStepResolver,
	) (models.SimulationStepResult, error)
}
