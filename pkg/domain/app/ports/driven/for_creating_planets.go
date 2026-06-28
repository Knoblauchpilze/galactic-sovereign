package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

type ForCreatingPlanets interface {
	Create(ctx context.Context, planet models.Planet) error
}
