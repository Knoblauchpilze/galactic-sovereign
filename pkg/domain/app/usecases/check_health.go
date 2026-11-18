package usecases

import (
	"context"

	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
)

type CheckHealthUseCase struct {
	checker drivenports.ForCheckingDatabaseConnection
}

func NewCheckHealthUseCase(checker drivenports.ForCheckingDatabaseConnection) *CheckHealthUseCase {
	return &CheckHealthUseCase{
		checker: checker,
	}
}

func (u *CheckHealthUseCase) Healthy(ctx context.Context) bool {
	err := u.checker.Ping(ctx)
	return err == nil
}
