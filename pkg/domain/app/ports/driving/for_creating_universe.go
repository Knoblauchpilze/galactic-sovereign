package driving

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type ForCreatingUniverse interface {
	Create(ctx context.Context, req request.UniverseCreationRequest) (models.Universe, error)
}
