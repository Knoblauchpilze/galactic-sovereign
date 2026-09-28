package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForCreatingFleet interface {
	Create(ctx context.Context, planet uuid.UUID, order models.FleetOrder) (models.Fleet, error)
}
