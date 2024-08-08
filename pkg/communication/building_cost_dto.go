package communication

import (
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingCostDtoResponse struct {
	Building uuid.UUID `json:"building"`
	Resource uuid.UUID `json:"resource"`
	Cost     int       `json:"cost"`
}

func ToBuildingCostDtoResponse(cost persistence.BuildingCost) BuildingCostDtoResponse {
	return BuildingCostDtoResponse{
		Building: cost.Building,
		Resource: cost.Resource,
		Cost:     cost.Cost,
	}
}
