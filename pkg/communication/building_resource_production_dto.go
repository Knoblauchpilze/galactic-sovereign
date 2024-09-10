package communication

import (
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingResourceProductionDtoResponse struct {
	Building uuid.UUID `json:"building"`
	Resource uuid.UUID `json:"resource"`
	Base     int       `json:"base"`
	Progress float64   `json:"progress"`
}

func ToBuildingResourceProductionDtoResponse(cost persistence.BuildingResourceProduction) BuildingResourceProductionDtoResponse {
	return BuildingResourceProductionDtoResponse{
		Building: cost.Building,
		Resource: cost.Resource,
		Base:     cost.Base,
		Progress: cost.Progress,
	}
}
