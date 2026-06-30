package drivenports

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type PlanetMutator func(*models.Planet, time.Time) error

type ForMutatingPlanets interface {
	Mutate(
		ctx context.Context,
		id uuid.UUID,
		mutator PlanetMutator,
	) (models.Planet, error)
}
