package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type SolarSystemRepository struct {
	conn db.Connection
}

func NewSolarSystemRepository(conn db.Connection) *SolarSystemRepository {
	return &SolarSystemRepository{
		conn: conn,
	}
}

func (r *SolarSystemRepository) GetSolarSystem(
	ctx context.Context,
	universe uuid.UUID,
	galaxy int,
	solarSystem int,
) (models.SolarSystem, error) {
	return models.SolarSystem{}, errors.ErrNotImplemented
}
