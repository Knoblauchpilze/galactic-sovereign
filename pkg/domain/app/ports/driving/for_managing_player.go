package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
)

type ForManagingPlayer interface {
	Create(ctx context.Context, req request.PlayerCreationRequest) (models.Player, error)
	Get(ctx context.Context, id uuid.UUID) (models.Player, error)
	ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]models.Player, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
