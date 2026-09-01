package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForFetchingBuilding interface {
	Get(ctx context.Context, id uuid.UUID) (models.Building, error)
}
