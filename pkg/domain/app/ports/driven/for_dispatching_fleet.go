package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type FleetCreator func(*models.Planet) (models.Fleet, error)

type ForDispatchingFleet interface {
	Dispatch(
		ctx context.Context,
		source uuid.UUID,
		mutator FleetCreator,
	) (models.Fleet, error)
}
