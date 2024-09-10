package communication

import (
	"encoding/json"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

type FullBuildingDtoResponse struct {
	BuildingDtoResponse

	Costs       []BuildingCostDtoResponse
	Productions []BuildingResourceProductionDtoResponse
}

func ToFullBuildingDtoResponse(building persistence.Building, costs []persistence.BuildingCost, productions []persistence.BuildingResourceProduction) FullBuildingDtoResponse {
	out := FullBuildingDtoResponse{
		BuildingDtoResponse: ToBuildingDtoResponse(building),
	}

	for _, cost := range costs {
		costDto := ToBuildingCostDtoResponse(cost)
		out.Costs = append(out.Costs, costDto)
	}

	for _, production := range productions {
		productionDto := ToBuildingResourceProductionDtoResponse(production)
		out.Productions = append(out.Productions, productionDto)
	}

	return out
}

func (dto FullBuildingDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		BuildingDtoResponse
		Costs       []BuildingCostDtoResponse               `json:"costs"`
		Productions []BuildingResourceProductionDtoResponse `json:"productions"`
	}{
		BuildingDtoResponse: dto.BuildingDtoResponse,
		Costs:               dto.Costs,
		Productions:         dto.Productions,
	}

	if out.Costs == nil {
		out.Costs = make([]BuildingCostDtoResponse, 0)
	}

	if out.Productions == nil {
		out.Productions = make([]BuildingResourceProductionDtoResponse, 0)
	}

	return json.Marshal(out)
}
