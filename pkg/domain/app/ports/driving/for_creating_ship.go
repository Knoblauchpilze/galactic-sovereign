package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type ForCreatingShip interface {
	Create(ctx context.Context, req request.ShipCreationRequest) error
}
