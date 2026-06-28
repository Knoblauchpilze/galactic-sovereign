package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type ForCreatingPlanet interface {
	Create(ctx context.Context, req request.PlanetCreationRequest) (models.Planet, error)
}
