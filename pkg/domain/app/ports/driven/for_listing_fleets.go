package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForListingFleets interface {
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Fleet, error)
}
