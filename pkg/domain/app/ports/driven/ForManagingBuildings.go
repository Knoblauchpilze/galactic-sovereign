package driven

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

type ForManagingBuildings interface {
	List(ctx context.Context) ([]models.Building, error)
}
