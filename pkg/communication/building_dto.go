package communication

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingDtoResponse struct {
	Id   uuid.UUID `json:"id" format:"uuid"`
	Name string    `json:"name"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
}

func ToBuildingDtoResponse(building persistence.Building) BuildingDtoResponse {
	return BuildingDtoResponse{
		Id:   building.Id,
		Name: building.Name,

		CreatedAt: building.CreatedAt,
	}
}
