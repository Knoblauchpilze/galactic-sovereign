package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
)

type universeUseCase struct {
	repo driven.ForManagingUniverses
}

func NewUniverseUseCase(repo driven.ForManagingUniverses) driving.ForCreatingUniverse {
	return &universeUseCase{
		repo: repo,
	}
}

func (u *universeUseCase) Create(ctx context.Context, req request.UniverseCreationRequest) (models.Universe, error) {
	universe := request.FromUniverseCreationRequest(req)

	err := u.repo.Create(ctx, universe)
	if err != nil {
		return models.Universe{}, err
	}

	return universe, nil
}
