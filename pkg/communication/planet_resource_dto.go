package communication

import (
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceDtoResponse struct {
	Planet   uuid.UUID `json:"planet"`
	Resource uuid.UUID `json:"resource"`
	Amount   float64   `json:"amount"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToPlanetResourceDtoResponse(resource persistence.PlanetResource) PlanetResourceDtoResponse {
	return PlanetResourceDtoResponse{
		Planet:   resource.Planet,
		Resource: resource.Resource,
		Amount:   resource.Amount,

		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
	}
}
