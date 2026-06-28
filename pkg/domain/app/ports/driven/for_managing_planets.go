package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForManagingPlanets interface {
	Get(ctx context.Context, id uuid.UUID) (models.Planet, error)
	GetByAction(ctx context.Context, action uuid.UUID) (models.Planet, error)
	List(ctx context.Context) ([]models.Planet, error)
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
