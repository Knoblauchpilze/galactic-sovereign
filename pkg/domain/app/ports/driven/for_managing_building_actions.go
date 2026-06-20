package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForManagingBuildingActions interface {
	Create(ctx context.Context, planet models.Planet) error
	Get(ctx context.Context, id uuid.UUID) (models.BuildingAction, error)
	Delete(ctx context.Context, planet models.Planet, action uuid.UUID) error
}
