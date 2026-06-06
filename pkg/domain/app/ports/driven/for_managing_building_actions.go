package drivenport

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForManagingBuildingActions interface {
	Create(ctx context.Context, action models.BuildingAction) error
	Get(ctx context.Context, id uuid.UUID) (models.BuildingAction, error)
	ListForPlanet(ctx context.Context, planet uuid.UUID) ([]models.BuildingAction, error)
	ListBeforeCompletionTime(ctx context.Context, planet uuid.UUID, until time.Time) ([]models.BuildingAction, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteForPlanet(ctx context.Context, planet uuid.UUID) error
	DeleteForPlayer(ctx context.Context, player uuid.UUID) error
}
