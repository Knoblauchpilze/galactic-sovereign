package communication

import (
	"encoding/json"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

type FullBuildingDtoResponse struct {
	BuildingDtoResponse

	Costs []BuildingCostDtoResponse
}

func ToFullBuildingDtoResponse(building persistence.Building, costs []persistence.BuildingCost, buildings []persistence.Building) FullBuildingDtoResponse {
	out := FullBuildingDtoResponse{
		BuildingDtoResponse: ToBuildingDtoResponse(building),
	}

	for _, cost := range costs {
		costDto := ToBuildingCostDtoResponse(cost)
		out.Costs = append(out.Costs, costDto)
	}

	return out
}

func (dto FullBuildingDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		BuildingDtoResponse
		Costs []BuildingCostDtoResponse `json:"costs"`
	}{
		BuildingDtoResponse: dto.BuildingDtoResponse,
		Costs:               dto.Costs,
	}

	if out.Costs == nil {
		out.Costs = make([]BuildingCostDtoResponse, 0)
	}

	return json.Marshal(out)
}
