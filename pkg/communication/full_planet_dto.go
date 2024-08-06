package communication

import (
	"encoding/json"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

type FullPlanetDtoResponse struct {
	PlanetDtoResponse

	Resources []PlanetResourceDtoResponse
	Buildings []PlanetBuildingDtoResponse
}

func ToFullPlanetDtoResponse(planet persistence.Planet, resources []persistence.PlanetResource, buildings []persistence.PlanetBuilding) FullPlanetDtoResponse {
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

	return out
}

func (dto FullPlanetDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		PlanetDtoResponse
		Resources []PlanetResourceDtoResponse `json:"resources"`
		Buildings []PlanetBuildingDtoResponse `json:"buildings"`
	}{
		PlanetDtoResponse: dto.PlanetDtoResponse,
		Resources:         dto.Resources,
		Buildings:         dto.Buildings,
	}
	if out.Resources == nil {
		out.Resources = make([]PlanetResourceDtoResponse, 0)
	}
	if out.Buildings == nil {
		out.Buildings = make([]PlanetBuildingDtoResponse, 0)
	}

	return json.Marshal(out)
}
