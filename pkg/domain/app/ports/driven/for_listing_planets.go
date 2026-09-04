package drivenports

import (
	"context"

	"github.com/google/uuid"
)

type ForListingPlanets interface {
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]uuid.UUID, error)
}
