package communication

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceDtoResponse struct {
	Planet   uuid.UUID `json:"planet" format:"uuid"`
	Resource uuid.UUID `json:"resource" format:"uuid"`
	Amount   float64   `json:"amount"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`
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
