package communication

import (
	"encoding/json"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

type FullPlanetDtoResponse struct {
	PlanetDtoResponse

	Resources []PlanetResourceDtoResponse
	Buildings []PlanetBuildingDtoResponse

	BuildingActions []BuildingActionDtoResponse
}

func ToFullPlanetDtoResponse(planet persistence.Planet, resources []persistence.PlanetResource, buildings []persistence.PlanetBuilding, buildingActions []persistence.BuildingAction) FullPlanetDtoResponse {
	out := FullPlanetDtoResponse{
		PlanetDtoResponse: ToPlanetDtoResponse(planet),
	}

	for _, resource := range resources {
		resourceDto := ToPlanetResourceDtoResponse(resource)
		out.Resources = append(out.Resources, resourceDto)
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
		Resources []PlanetResourceDtoResponse `json:"resources"`
		Buildings []PlanetBuildingDtoResponse `json:"buildings"`

		BuildingActions []BuildingActionDtoResponse `json:"buildingActions"`
	}{
		PlanetDtoResponse: dto.PlanetDtoResponse,
		Resources:         dto.Resources,
		Buildings:         dto.Buildings,

		BuildingActions: dto.BuildingActions,
	}
	if out.Resources == nil {
		out.Resources = make([]PlanetResourceDtoResponse, 0)
	}
	if out.Buildings == nil {
		out.Buildings = make([]PlanetBuildingDtoResponse, 0)
	}
	if out.BuildingActions == nil {
		out.BuildingActions = make([]BuildingActionDtoResponse, 0)
	}

	return json.Marshal(out)
}
