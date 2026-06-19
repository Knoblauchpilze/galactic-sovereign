package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
)

type ForManagingPlanet interface {
	Create(ctx context.Context, req request.PlanetCreationRequest) (models.Planet, error)
	Get(ctx context.Context, id uuid.UUID) (models.Planet, error)
	List(ctx context.Context) ([]models.Planet, error)
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Planet, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
