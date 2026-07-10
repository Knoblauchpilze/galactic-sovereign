package drivenports

import (
	"context"

	"github.com/google/uuid"
)

type ForManagingPlanets interface {
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
