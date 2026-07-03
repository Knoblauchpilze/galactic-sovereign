package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type PlanetMutator func(*models.Planet) error

type ForMutatingPlanet interface {
	Mutate(
		ctx context.Context,
		id uuid.UUID,
		mutator PlanetMutator,
	) (models.Planet, error)
}
