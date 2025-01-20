package communication

import (
	"encoding/json"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type FullUniverseDtoResponse struct {
	UniverseDtoResponse

	Resources []ResourceDtoResponse
	Buildings []FullBuildingDtoResponse
}

func ToFullUniverseDtoResponse(universe persistence.Universe,
	resources []persistence.Resource,
	buildings []persistence.Building,
	buildingCosts map[uuid.UUID][]persistence.BuildingCost,
	buildingProductions map[uuid.UUID][]persistence.BuildingResourceProduction) FullUniverseDtoResponse {
	out := FullUniverseDtoResponse{
		UniverseDtoResponse: ToUniverseDtoResponse(universe),
	}

	for _, resource := range resources {
		resourceDto := ToResourceDtoResponse(resource)
		out.Resources = append(out.Resources, resourceDto)
	}

	for _, building := range buildings {
		buildingDto := ToFullBuildingDtoResponse(building, buildingCosts[building.Id], buildingProductions[building.Id])
		out.Buildings = append(out.Buildings, buildingDto)
	}

	return out
}

func (dto FullUniverseDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		UniverseDtoResponse
		Resources []ResourceDtoResponse     `json:"resources"`
		Buildings []FullBuildingDtoResponse `json:"buildings"`
	}{
		UniverseDtoResponse: dto.UniverseDtoResponse,
		Resources:           dto.Resources,
		Buildings:           dto.Buildings,
	}

	if out.Resources == nil {
		out.Resources = make([]ResourceDtoResponse, 0)
	}
	if out.Buildings == nil {
		out.Buildings = make([]FullBuildingDtoResponse, 0)
	}

	return json.Marshal(out)
}
