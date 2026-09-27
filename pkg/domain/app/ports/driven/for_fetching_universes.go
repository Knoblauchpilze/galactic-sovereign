package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForFetchingUniverses interface {
	Get(ctx context.Context, id uuid.UUID) (models.Universe, error)
	GetByPlanetId(ctx context.Context, planet uuid.UUID) (models.Universe, error)
	List(ctx context.Context) ([]models.Universe, error)
}
