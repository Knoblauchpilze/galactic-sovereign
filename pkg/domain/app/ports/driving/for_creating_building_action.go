package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type ForCreatingBuildingAction interface {
	Create(ctx context.Context, req request.BuildingActionCreationRequest) (models.BuildingAction, error)
}
