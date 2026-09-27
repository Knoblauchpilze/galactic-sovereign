package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForManagingUniverses interface {
	Create(ctx context.Context, universe models.Universe) error
	Delete(ctx context.Context, id uuid.UUID) error
}
