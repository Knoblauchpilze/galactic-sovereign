package communication

import (
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingDtoResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
}

func ToBuildingDtoResponse(building persistence.Building) BuildingDtoResponse {
	return BuildingDtoResponse{
		Id:   building.Id,
		Name: building.Name,

		CreatedAt: building.CreatedAt,
	}
}
