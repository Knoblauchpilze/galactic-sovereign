package drivingports

import (
	"context"

	"github.com/google/uuid"
)

type ForDeletingBuildingAction interface {
	DeleteForPlanet(ctx context.Context, planet uuid.UUID) error
}
