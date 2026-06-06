package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
)

type universeUseCase struct {
	repo drivenports.ForManagingUniverses
}

func NewUniverseUseCase(repo drivenports.ForManagingUniverses) drivingports.ForManagingUniverse {
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

func (u *universeUseCase) Get(ctx context.Context, id uuid.UUID) (models.Universe, error) {
	return u.repo.Get(ctx, id)
}

func (u *universeUseCase) List(ctx context.Context) ([]models.Universe, error) {
	return u.repo.List(ctx)
}

func (u *universeUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
