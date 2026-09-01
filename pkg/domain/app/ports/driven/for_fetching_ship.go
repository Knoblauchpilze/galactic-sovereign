package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForFetchingShip interface {
	Get(ctx context.Context, id uuid.UUID) (models.Ship, error)
}
