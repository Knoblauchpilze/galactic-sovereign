package communication

import (
	"encoding/json"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

type FullUniverseDtoResponse struct {
	UniverseDtoResponse

	Resources []ResourceDtoResponse
}

func ToFullUniverseDtoResponse(universe persistence.Universe, resources []persistence.Resource) FullUniverseDtoResponse {
	out := FullUniverseDtoResponse{
		UniverseDtoResponse: ToUniverseDtoResponse(universe),
	}

	for _, resource := range resources {
		resourceDto := ToResourceDtoResponse(resource)
		out.Resources = append(out.Resources, resourceDto)
	}

	return out
}

func (dto FullUniverseDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		UniverseDtoResponse
		Resources []ResourceDtoResponse `json:"resources"`
	}{
		UniverseDtoResponse: dto.UniverseDtoResponse,
		Resources:           dto.Resources,
	}
	if out.Resources == nil {
		out.Resources = make([]ResourceDtoResponse, 0)
	}

	return json.Marshal(out)
}
