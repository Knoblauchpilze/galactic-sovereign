package communication

import (
	"encoding/json"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
)

type FullPlanetDtoResponse struct {
	PlanetDtoResponse

	Resources   []PlanetResourceDtoResponse
	Productions []PlanetResourceProductionDtoResponse
	Storages    []PlanetResourceStorageDtoResponse
	Buildings   []PlanetBuildingDtoResponse

	BuildingActions []BuildingActionDtoResponse
}

func ToFullPlanetDtoResponse(planet persistence.Planet,
	resources []persistence.PlanetResource,
	productions []persistence.PlanetResourceProduction,
	storages []persistence.PlanetResourceStorage,
	buildings []persistence.PlanetBuilding,
	buildingActions []persistence.BuildingAction) FullPlanetDtoResponse {
	out := FullPlanetDtoResponse{
		PlanetDtoResponse: ToPlanetDtoResponse(planet),
	}

	for _, resource := range resources {
		resourceDto := ToPlanetResourceDtoResponse(resource)
		out.Resources = append(out.Resources, resourceDto)
	}

	for _, production := range productions {
		productionDto := ToPlanetResourceProductionDtoResponse(production)
		out.Productions = append(out.Productions, productionDto)
	}

	for _, storage := range storages {
		storageDto := ToPlanetResourceStorageDtoResponse(storage)
		out.Storages = append(out.Storages, storageDto)
	}

	for _, building := range buildings {
		buildingDto := ToPlanetBuildingDtoResponse(building)
		out.Buildings = append(out.Buildings, buildingDto)
	}

	for _, action := range buildingActions {
		actionDto := ToBuildingActionDtoResponse(action)
		out.BuildingActions = append(out.BuildingActions, actionDto)
	}

	return out
}

func (dto FullPlanetDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		PlanetDtoResponse
		Resources   []PlanetResourceDtoResponse           `json:"resources"`
		Productions []PlanetResourceProductionDtoResponse `json:"productions"`
		Storages    []PlanetResourceStorageDtoResponse    `json:"storages"`
		Buildings   []PlanetBuildingDtoResponse           `json:"buildings"`

		BuildingActions []BuildingActionDtoResponse `json:"buildingActions"`
	}{
		PlanetDtoResponse: dto.PlanetDtoResponse,
		Resources:         dto.Resources,
		Productions:       dto.Productions,
		Storages:          dto.Storages,
		Buildings:         dto.Buildings,

		BuildingActions: dto.BuildingActions,
	}
	if out.Resources == nil {
		out.Resources = make([]PlanetResourceDtoResponse, 0)
	}
	if out.Productions == nil {
		out.Productions = make([]PlanetResourceProductionDtoResponse, 0)
	}
	if out.Storages == nil {
		out.Storages = make([]PlanetResourceStorageDtoResponse, 0)
	}
	if out.Buildings == nil {
		out.Buildings = make([]PlanetBuildingDtoResponse, 0)
	}
	if out.BuildingActions == nil {
		out.BuildingActions = make([]BuildingActionDtoResponse, 0)
	}

	return json.Marshal(out)
}
