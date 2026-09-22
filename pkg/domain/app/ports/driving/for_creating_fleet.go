package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type ForCreatingFleet interface {
	Create(ctx context.Context, req request.FleetCreationRequest) (models.Fleet, error)
}
