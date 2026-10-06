package drivingports

import (
	"context"

	"github.com/google/uuid"
)

type ForDeletingPlanet interface {
	Delete(ctx context.Context, id uuid.UUID) error
}
