package communication

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetBuildingDtoResponse struct {
	Planet   uuid.UUID `json:"planet"`
	Building uuid.UUID `json:"building"`
	Level    int       `json:"level"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToPlanetBuildingDtoResponse(building persistence.PlanetBuilding) PlanetBuildingDtoResponse {
	return PlanetBuildingDtoResponse{
		Planet:   building.Planet,
		Building: building.Building,
		Level:    building.Level,

		CreatedAt: building.CreatedAt,
		UpdatedAt: building.UpdatedAt,
	}
}
