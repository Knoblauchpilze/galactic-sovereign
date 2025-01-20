package communication

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingCostDtoResponse struct {
	Building uuid.UUID `json:"building"`
	Resource uuid.UUID `json:"resource"`
	Cost     int       `json:"cost"`
	Progress float64   `json:"progress"`
}

func ToBuildingCostDtoResponse(cost persistence.BuildingCost) BuildingCostDtoResponse {
	return BuildingCostDtoResponse{
		Building: cost.Building,
		Resource: cost.Resource,
		Cost:     cost.Cost,
		Progress: cost.Progress,
	}
}
