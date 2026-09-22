package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

type CreateFleetUseCase struct {
}

func NewCreateFleetUseCase() *CreateFleetUseCase {
	return &CreateFleetUseCase{}
}

func (b *CreateFleetUseCase) Create(
	ctx context.Context,
	req request.FleetCreationRequest,
) (models.Fleet, error) {
	return models.Fleet{}, nil

}
