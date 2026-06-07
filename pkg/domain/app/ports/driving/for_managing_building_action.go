package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
)

type ForManagingBuildingAction interface {
	Create(ctx context.Context, req request.BuildingActionCreationRequest) (models.BuildingAction, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
