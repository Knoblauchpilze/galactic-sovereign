package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

func ToUniverseCreationRequest(dto dtos.UniverseDtoRequest) request.UniverseCreationRequest {
	return request.UniverseCreationRequest{
		Name: dto.Name,
		Topology: request.TopologyRequest{
			Galaxies:     dto.Topology.Galaxies,
			SolarSystems: dto.Topology.SolarSystems,
			Orbits:       dto.Topology.Orbits,
		},
	}
}

func ToUniverseResponse(universe models.Universe) dtos.UniverseDtoResponse {
	return dtos.UniverseDtoResponse{
		Id:        universe.Id,
		Name:      universe.Name,
		CreatedAt: universe.CreatedAt,
		Topology:  toTopologyResponse(universe.Topology),
		Resources: toResourcesResponse(universe.Resources),
		Buildings: toBuildingsResponse(universe.Buildings),
		Ships:     toShipsResponse(universe.Ships),
	}
}

func ToUniversesResponse(universes []models.Universe) []dtos.UniverseDtoResponse {
	out := make([]dtos.UniverseDtoResponse, 0, len(universes))

	for _, u := range universes {
		dto := ToUniverseResponse(u)
		out = append(out, dto)
	}

	return out
}

func toTopologyResponse(topology models.UniverseTopology) dtos.TopologyDtoResponse {
	return dtos.TopologyDtoResponse{
		Galaxies:     topology.Galaxies,
		SolarSystems: topology.SolarSystems,
		Orbits:       topology.Orbits,
	}
}

func toResourceResponse(
	resource models.Resource,
) dtos.ResourceDtoResponse {
	return dtos.ResourceDtoResponse{
		Id:                            resource.Id,
		Name:                          resource.Name,
		StartAmount:                   resource.StartAmount,
		StartProduction:               resource.StartProduction,
		StartStorage:                  resource.StartStorage,
		BuildingBuildTimeHoursPerUnit: resource.BuildingBuildTimeHoursPerUnit,
		ShipyardBuildTimeHoursPerUnit: resource.ShipyardBuildTimeHoursPerUnit,
		CreatedAt:                     resource.CreatedAt,
	}
}

func toResourcesResponse(
	resources []models.Resource,
) []dtos.ResourceDtoResponse {
	out := make([]dtos.ResourceDtoResponse, 0, len(resources))

	for _, r := range resources {
		dto := toResourceResponse(r)
		out = append(out, dto)
	}

	return out
}

func toBuildingResponse(
	building models.Building,
) dtos.BuildingDtoResponse {
	dto := dtos.BuildingDtoResponse{
		Id:          building.Id,
		Name:        building.Name,
		CreatedAt:   building.CreatedAt,
		Costs:       toBuildingCostsResponse(building.Costs),
		Productions: toBuildingProductionsResponse(building.Productions),
		Storages:    toBuildingStoragesResponse(building.Storages),
	}

	if building.ShipSpeedup != nil {
		speedup := toBuildingShipSpeedupResponse(*building.ShipSpeedup)
		dto.ShipSpeedup = &speedup
	}

	return dto
}

func toBuildingsResponse(
	buildings []models.Building,
) []dtos.BuildingDtoResponse {
	out := make([]dtos.BuildingDtoResponse, 0, len(buildings))

	for _, b := range buildings {
		dto := toBuildingResponse(b)
		out = append(out, dto)
	}

	return out
}

func toBuildingCostResponse(
	cost models.BuildingCost,
) dtos.BuildingCostDtoResponse {
	return dtos.BuildingCostDtoResponse{
		Resource: cost.Resource,
		Cost:     cost.Cost,
		Progress: cost.Progress,
	}
}

func toBuildingCostsResponse(
	costs []models.BuildingCost,
) []dtos.BuildingCostDtoResponse {
	out := make([]dtos.BuildingCostDtoResponse, 0, len(costs))

	for _, c := range costs {
		dto := toBuildingCostResponse(c)
		out = append(out, dto)
	}

	return out
}

func toBuildingProductionResponse(
	production models.BuildingResourceProduction,
) dtos.BuildingResourceProductionDtoResponse {
	return dtos.BuildingResourceProductionDtoResponse{
		Resource: production.Resource,
		Base:     production.Base,
		Progress: production.Progress,
	}
}

func toBuildingProductionsResponse(
	productions []models.BuildingResourceProduction,
) []dtos.BuildingResourceProductionDtoResponse {
	out := make([]dtos.BuildingResourceProductionDtoResponse, 0, len(productions))

	for _, p := range productions {
		dto := toBuildingProductionResponse(p)
		out = append(out, dto)
	}

	return out
}

func toBuildingStorageResponse(
	storage models.BuildingResourceStorage,
) dtos.BuildingResourceStorageDtoResponse {
	return dtos.BuildingResourceStorageDtoResponse{
		Resource: storage.Resource,
		Base:     storage.Base,
		Scale:    storage.Scale,
		Progress: storage.Progress,
	}
}

func toBuildingStoragesResponse(
	storages []models.BuildingResourceStorage,
) []dtos.BuildingResourceStorageDtoResponse {
	out := make([]dtos.BuildingResourceStorageDtoResponse, 0, len(storages))

	for _, s := range storages {
		dto := toBuildingStorageResponse(s)
		out = append(out, dto)
	}

	return out
}

func toShipResponse(
	ship models.Ship,
) dtos.ShipDtoResponse {
	return dtos.ShipDtoResponse{
		Id:                   ship.Id,
		Name:                 ship.Name,
		CreatedAt:            ship.CreatedAt,
		Costs:                toShipCostsResponse(ship.Costs),
		BuildingRequirements: toShipBuildingRequirementsResponse(ship.BuildingRequirements),
	}
}

func toShipsResponse(
	ships []models.Ship,
) []dtos.ShipDtoResponse {
	out := make([]dtos.ShipDtoResponse, 0, len(ships))

	for _, s := range ships {
		dto := toShipResponse(s)
		out = append(out, dto)
	}

	return out
}

func toShipCostResponse(
	cost models.ShipCost,
) dtos.ShipCostDtoResponse {
	return dtos.ShipCostDtoResponse{
		Resource: cost.Resource,
		Cost:     cost.Cost,
	}
}

func toShipCostsResponse(
	costs []models.ShipCost,
) []dtos.ShipCostDtoResponse {
	out := make([]dtos.ShipCostDtoResponse, 0, len(costs))

	for _, c := range costs {
		dto := toShipCostResponse(c)
		out = append(out, dto)
	}

	return out
}

func toShipBuildingRequirementResponse(
	requirement models.ShipBuildingRequirement,
) dtos.ShipBuildingRequirementDtoResponse {
	return dtos.ShipBuildingRequirementDtoResponse{
		Building: requirement.Building,
		Level:    requirement.Level,
	}
}

func toShipBuildingRequirementsResponse(
	requirements []models.ShipBuildingRequirement,
) []dtos.ShipBuildingRequirementDtoResponse {
	out := make([]dtos.ShipBuildingRequirementDtoResponse, 0, len(requirements))

	for _, r := range requirements {
		dto := toShipBuildingRequirementResponse(r)
		out = append(out, dto)
	}

	return out
}

func toBuildingShipSpeedupResponse(
	shipSpeedup models.BuildingShipSpeedup,
) dtos.BuildingShipSpeedupDtoResponse {
	return dtos.BuildingShipSpeedupDtoResponse{
		Scaling:  shipSpeedup.Scaling,
		Base:     shipSpeedup.Base,
		Progress: shipSpeedup.Progress,
	}
}
