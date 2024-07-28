package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ResourceDtoResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
}

func ToResourceDtoResponse(resource persistence.Resource) ResourceDtoResponse {
	return ResourceDtoResponse{
		Id:   resource.Id,
		Name: resource.Name,

		CreatedAt: resource.CreatedAt,
	}
}
