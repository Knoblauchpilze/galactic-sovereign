package persistence

import (
	"time"

	"github.com/google/uuid"
)

type BuildingAction struct {
	Id           uuid.UUID
	Planet       uuid.UUID
	Building     uuid.UUID
	CurrentLevel int
	DesiredLevel int
	CreatedAt    time.Time
	CompletedAt  time.Time
}

func ToPlanetBuilding(action BuildingAction, building PlanetBuilding) PlanetBuilding {
	return PlanetBuilding{
		Planet:   action.Planet,
		Building: action.Building,
		Level:    action.DesiredLevel,

		CreatedAt: building.CreatedAt,
		UpdatedAt: action.CompletedAt,

		Version: building.Version,
	}
}
