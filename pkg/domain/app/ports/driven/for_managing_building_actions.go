package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

// TODO: This can be removed and replaced by a mutation of the planet
type ForManagingBuildingActions interface {
	Create(ctx context.Context, planet models.Planet) error
	Delete(ctx context.Context, planet models.Planet) error
}
