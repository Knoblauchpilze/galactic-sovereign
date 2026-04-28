package communication

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetBuildingDtoResponse struct {
	Planet   uuid.UUID `json:"planet" format:"uuid"`
	Building uuid.UUID `json:"building" format:"uuid"`
	Level    int       `json:"level"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`
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
