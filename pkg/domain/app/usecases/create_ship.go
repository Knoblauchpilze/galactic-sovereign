package usecases

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type CreateShipUseCase struct{}

func NewCreateShipUseCase() *CreateShipUseCase {
	return &CreateShipUseCase{}
}

func (b *CreateShipUseCase) Create(
	ctx context.Context,
	req request.ShipCreationRequest,
) error {
	return errors.ErrNotImplemented
}
