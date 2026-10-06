package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForFetchingPlanets interface {
	Get(ctx context.Context, id uuid.UUID) (models.Planet, error)
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error)
}
