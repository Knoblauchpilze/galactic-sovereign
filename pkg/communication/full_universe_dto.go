package communication

import (
	"encoding/json"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

type FullUniverseDtoResponse struct {
	UniverseDtoResponse

	Resources []ResourceDtoResponse
	Buildings []BuildingDtoResponse
}

func ToFullUniverseDtoResponse(universe persistence.Universe, resources []persistence.Resource, buildings []persistence.Building) FullUniverseDtoResponse {
	out := FullUniverseDtoResponse{
		UniverseDtoResponse: ToUniverseDtoResponse(universe),
	}

	for _, resource := range resources {
		resourceDto := ToResourceDtoResponse(resource)
		out.Resources = append(out.Resources, resourceDto)
	}

	for _, building := range buildings {
		buildingDto := ToBuildingDtoResponse(building)
		out.Buildings = append(out.Buildings, buildingDto)
	}

	return out
}

func (dto FullUniverseDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		UniverseDtoResponse
		Resources []ResourceDtoResponse `json:"resources"`
		Buildings []BuildingDtoResponse `json:"buildings"`
	}{
		UniverseDtoResponse: dto.UniverseDtoResponse,
		Resources:           dto.Resources,
		Buildings:           dto.Buildings,
	}

	if out.Resources == nil {
		out.Resources = make([]ResourceDtoResponse, 0)
	}
	if out.Buildings == nil {
		out.Buildings = make([]BuildingDtoResponse, 0)
	}

	return json.Marshal(out)
}
