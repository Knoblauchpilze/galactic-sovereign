package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UniverseDtoRequest struct {
	Name string `json:"name" form:"name"`
}

type UniverseDtoResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
}

func FromUniverseDtoRequest(universe UniverseDtoRequest) persistence.Universe {
	t := time.Now()
	return persistence.Universe{
		Id:   uuid.New(),
		Name: universe.Name,

		CreatedAt: t,
		UpdatedAt: t,
	}
}

func ToUniverseDtoResponse(universe persistence.Universe) UniverseDtoResponse {
	return UniverseDtoResponse{
		Id:   universe.Id,
		Name: universe.Name,

		CreatedAt: universe.CreatedAt,
	}
}
