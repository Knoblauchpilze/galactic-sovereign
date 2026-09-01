package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type ForCreatingShipAction interface {
	Create(ctx context.Context, req request.ShipActionCreationRequest) (models.ShipAction, error)
}
