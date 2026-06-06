package drivenport

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

type ForListingBuildings interface {
	List(ctx context.Context) ([]models.Building, error)
}
