package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

// PlanetMutator a function called by the mutator with the latest version of a
// planet loaded from the database. The function is allowed to mutate the planet
// and the mutator will persist the updated version to the database. In case
// any error happens during the process, the data is not persisted and this
// operation results in a no-op. The error is returned to the caller.
// The mutator is allowed to return a boolean: this boolean indicates whether
// the planet needs to be deleted.
type PlanetMutator func(*models.Planet) (bool, error)

type ForMutatingPlanet interface {
	Mutate(
		ctx context.Context,
		id uuid.UUID,
		mutator PlanetMutator,
	) (models.PlanetMutationResult, error)
}
