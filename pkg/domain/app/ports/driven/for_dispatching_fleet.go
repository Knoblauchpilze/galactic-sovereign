package drivenports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type FleetCreator func(*models.Planet, request.FleetCreationRequest) (models.Fleet, error)

type ForDispatchingFleet interface {
	Dispatch(
		ctx context.Context,
		req request.FleetCreationRequest,
		mutator FleetCreator,
	) (models.Fleet, error)
}
