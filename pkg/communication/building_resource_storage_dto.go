package communication

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingResourceStorageDtoResponse struct {
	Building uuid.UUID `json:"building"`
	Resource uuid.UUID `json:"resource"`
	Base     int       `json:"base"`
	Scale    float64   `json:"scale"`
	Progress float64   `json:"progress"`
}

func ToBuildingResourceStorageDtoResponse(storage persistence.BuildingResourceStorage) BuildingResourceStorageDtoResponse {
	return BuildingResourceStorageDtoResponse{
		Building: storage.Building,
		Resource: storage.Resource,
		Base:     storage.Base,
		Scale:    storage.Scale,
		Progress: storage.Progress,
	}
}
