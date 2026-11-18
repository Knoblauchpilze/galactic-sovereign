package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type FetchFleetsUseCase struct {
}

func NewFetchFleetsUseCase() *FetchFleetsUseCase {
	return &FetchFleetsUseCase{}
}

func (u *FetchFleetsUseCase) ListForPlayer(ctx context.Context, player uuid.UUID) ([]models.Fleet, error) {
	return []models.Fleet{}, nil
}
