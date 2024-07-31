package communication

import (
	"encoding/json"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

type FullPlanetDtoResponse struct {
	PlanetDtoResponse

	Resources []PlanetResourceDtoResponse `json:"resources"`
}

func ToFullPlanetDtoResponse(planet persistence.Planet, resources []persistence.PlanetResource) FullPlanetDtoResponse {
	out := FullPlanetDtoResponse{
		PlanetDtoResponse: ToPlanetDtoResponse(planet),
	}

	for _, resource := range resources {
		resourceDto := ToPlanetResourceDtoResponse(resource)
		out.Resources = append(out.Resources, resourceDto)
	}

	return out
}

func (dto FullPlanetDtoResponse) MarshalJSON() ([]byte, error) {
	out := struct {
		PlanetDtoResponse
		Resources []PlanetResourceDtoResponse `json:"resources"`
	}{
		PlanetDtoResponse: dto.PlanetDtoResponse,
		Resources:         dto.Resources,
	}
	if out.Resources == nil {
		out.Resources = make([]PlanetResourceDtoResponse, 0)
	}

	return json.Marshal(out)
}
