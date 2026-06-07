package usecases

import (
	"context"

	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
)

type checkHealthUseCase struct {
	checker drivenports.ForCheckingDatabaseConnection
}

func NewCheckHealthUseCase(checker drivenports.ForCheckingDatabaseConnection) drivingports.ForCheckingServiceHealth {
	return &checkHealthUseCase{
		checker: checker,
	}
}

func (c *checkHealthUseCase) Healthy(ctx context.Context) bool {
	err := c.checker.Ping(ctx)
	return err == nil
}
