package driving

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
)

type ForManagingUniverse interface {
	Create(ctx context.Context, req request.UniverseCreationRequest) (models.Universe, error)
	List(ctx context.Context) ([]models.Universe, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
