package drivingports

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type ForFetchingSolarSystem interface {
	GetSolarSystem(
		ctx context.Context,
		universe uuid.UUID,
		galaxy int,
		solarSystem int,
	) (models.SolarSystem, error)
}
