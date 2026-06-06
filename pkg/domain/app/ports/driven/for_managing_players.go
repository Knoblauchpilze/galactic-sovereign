package drivenport

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForManagingPlayers interface {
	Create(ctx context.Context, player models.Player) error
	Get(ctx context.Context, id uuid.UUID) (models.Player, error)
	List(ctx context.Context) ([]models.Player, error)
	ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]models.Player, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
