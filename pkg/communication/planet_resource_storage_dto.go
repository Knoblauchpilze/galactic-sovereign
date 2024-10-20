package communication

import (
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceStorageDtoResponse struct {
	Planet   uuid.UUID `json:"planet"`
	Resource uuid.UUID `json:"resource"`
	Storage  int       `json:"storage"`
}

func ToPlanetResourceStorageDtoResponse(storage persistence.PlanetResourceStorage) PlanetResourceStorageDtoResponse {
	return PlanetResourceStorageDtoResponse{
		Planet:   storage.Planet,
		Resource: storage.Resource,
		Storage:  storage.Storage,
	}
}
