package communication

import (
	"encoding/json"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
)

type FullBuildingDtoResponse struct {
	BuildingDtoResponse

	Costs       []BuildingCostDtoResponse
	Productions []BuildingResourceProductionDtoResponse
	Storages    []BuildingResourceStorageDtoResponse
}

func ToFullBuildingDtoResponse(
	building persistence.Building,
	costs []persistence.BuildingCost,
	productions []persistence.BuildingResourceProduction,
	storages []persistence.BuildingResourceStorage) FullBuildingDtoResponse {
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

	for _, storage := range storages {
		storageDto := ToBuildingResourceStorageDtoResponse(storage)
		out.Storages = append(out.Storages, storageDto)
	}

	return out
}

func (dto FullBuildingDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		BuildingDtoResponse
		Costs       []BuildingCostDtoResponse               `json:"costs"`
		Productions []BuildingResourceProductionDtoResponse `json:"productions"`
		Storages    []BuildingResourceStorageDtoResponse    `json:"storages"`
	}{
		BuildingDtoResponse: dto.BuildingDtoResponse,
		Costs:               dto.Costs,
		Productions:         dto.Productions,
		Storages:            dto.Storages,
	}

	if out.Costs == nil {
		out.Costs = make([]BuildingCostDtoResponse, 0)
	}

	if out.Productions == nil {
		out.Productions = make([]BuildingResourceProductionDtoResponse, 0)
	}

	if out.Storages == nil {
		out.Storages = make([]BuildingResourceStorageDtoResponse, 0)
	}

	return json.Marshal(out)
}
