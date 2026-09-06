package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbBuilding struct {
	Id                  uuid.UUID
	Name                string
	CreatedAt           time.Time
	ShipSpeedupScaling  *DbScalingMode
	ShipSpeedupBase     *float64
	ShipSpeedupProgress *float64
}

func (b DbBuilding) ToDomain() models.Building {
	return models.Building{
		Id:          b.Id,
		Name:        b.Name,
		CreatedAt:   b.CreatedAt,
		ShipSpeedup: b.tryMapSpeedup(),
	}
}

func (b DbBuilding) tryMapSpeedup() *models.BuildingShipSpeedup {
	// Note: it should not be possible to have only one nil value
	// considering the database schema. Therefore any partial read
	// will be considered invalid.
	if b.ShipSpeedupScaling == nil {
		return nil
	}
	if b.ShipSpeedupBase == nil {
		return nil
	}
	if b.ShipSpeedupProgress == nil {
		return nil
	}

	return &models.BuildingShipSpeedup{
		Scaling:  b.ShipSpeedupScaling.ToDomain(),
		Base:     *b.ShipSpeedupBase,
		Progress: *b.ShipSpeedupProgress,
	}
}
