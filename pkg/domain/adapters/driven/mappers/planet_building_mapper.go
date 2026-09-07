package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbPlanetBuilding struct {
	Building            uuid.UUID
	Level               int
	ShipSpeedupScaling  *DbScalingMode
	ShipSpeedupBase     *float64
	ShipSpeedupProgress *float64
}

func (b DbPlanetBuilding) ToDomain() models.PlanetBuilding {
	return models.PlanetBuilding{
		Building:    b.Building,
		Level:       b.Level,
		ShipSpeedup: b.tryMapSpeedup(),
	}
}

func (b DbPlanetBuilding) tryMapSpeedup() *models.BuildingShipSpeedup {
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
