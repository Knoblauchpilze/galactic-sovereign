package communication

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type ResourceDtoResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// TODO: Should include start_production and other props

	CreatedAt time.Time `json:"createdAt"`
}

func ToResourceDtoResponse(resource persistence.Resource) ResourceDtoResponse {
	return ResourceDtoResponse{
		Id:   resource.Id,
		Name: resource.Name,

		CreatedAt: resource.CreatedAt,
	}
}
