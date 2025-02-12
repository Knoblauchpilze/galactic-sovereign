package communication

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingResourceProductionDtoResponse struct {
	Building uuid.UUID `json:"building"`
	Resource uuid.UUID `json:"resource"`
	Base     int       `json:"base"`
	Progress float64   `json:"progress"`
}

func ToBuildingResourceProductionDtoResponse(prod persistence.BuildingResourceProduction) BuildingResourceProductionDtoResponse {
	return BuildingResourceProductionDtoResponse{
		Building: prod.Building,
		Resource: prod.Resource,
		Base:     prod.Base,
		Progress: prod.Progress,
	}
}
