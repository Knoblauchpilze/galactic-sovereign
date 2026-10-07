package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

// PlanetDeleter a function called by the deleter with the latest version of a
// planet loaded from the database. The function can perform some checks on the
// planet before authorizing the deletion.
// In case an error is returned, the deletion will be halted and the adapter is
// expected to return the error from the deleter. If the closure does not return
// an error, the adapter should proceed with the planet's deletion.
type PlanetDeleter func(*models.Planet) error

type ForDeletingPlanet interface {
	Delete(
		ctx context.Context,
		id uuid.UUID,
		deleter PlanetDeleter,
	) error
}
