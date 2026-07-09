package drivingports

import (
	"context"

	"github.com/google/uuid"
)

type ForManagingBuildingAction interface {
	Delete(ctx context.Context, id uuid.UUID) error
}
