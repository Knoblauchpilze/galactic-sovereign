package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbUsedSlot struct {
	Planet      uuid.UUID
	Galaxy      int
	SolarSystem int
	Position    int
}

func (s DbUsedSlot) ToDomain() (uuid.UUID, models.Coordinate) {
	return s.Planet, models.Coordinate{
		Galaxy:      s.Galaxy,
		SolarSystem: s.SolarSystem,
		Position:    s.Position,
	}
}
