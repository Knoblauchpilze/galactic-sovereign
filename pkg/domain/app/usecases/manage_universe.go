package usecases

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

type UniverseUseCase struct {
	fetchRepo drivenports.ForFetchingUniverses
	repo      drivenports.ForManagingUniverses
}

func NewUniverseUseCase(
	fetchRepo drivenports.ForFetchingUniverses,
	repo drivenports.ForManagingUniverses) *UniverseUseCase {
	return &UniverseUseCase{
		fetchRepo: fetchRepo,
		repo:      repo,
	}
}

func (u *UniverseUseCase) Create(ctx context.Context, req request.UniverseCreationRequest) (models.Universe, error) {
	universe := request.FromUniverseCreationRequest(req)

	err := u.repo.Create(ctx, universe)
	if err != nil {
		return models.Universe{}, err
	}

	return universe, nil
}

func (u *UniverseUseCase) Get(ctx context.Context, id uuid.UUID) (models.Universe, error) {
	return u.fetchRepo.Get(ctx, id)
}

func (u *UniverseUseCase) List(ctx context.Context) ([]models.Universe, error) {
	return u.fetchRepo.List(ctx)
}

func (u *UniverseUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
