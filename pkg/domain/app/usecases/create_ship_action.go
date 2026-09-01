package usecases

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
)

type CreateShipActionUseCase struct {
	shipRepo      drivenports.ForFetchingShip
	planetMutator drivenports.ForMutatingPlanet
	clock         drivenports.ForFetchingTime
}

func NewCreateShipActionUseCase(
	shipRepo drivenports.ForFetchingShip,
	planetMutator drivenports.ForMutatingPlanet,
	clock drivenports.ForFetchingTime,
) *CreateShipActionUseCase {
	return &CreateShipActionUseCase{
		shipRepo:      shipRepo,
		planetMutator: planetMutator,
		clock:         clock,
	}
}

func (b *CreateShipActionUseCase) Create(
	ctx context.Context,
	req request.ShipActionCreationRequest,
) error {
	return errors.ErrNotImplemented
}
