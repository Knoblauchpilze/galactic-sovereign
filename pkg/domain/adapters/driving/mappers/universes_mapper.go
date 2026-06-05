package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

func ToUniverseCreationRequest(dto dtos.UniverseDtoRequest) request.UniverseCreationRequest {
	return request.UniverseCreationRequest{
		Name: dto.Name,
	}
}

func ToUniverseResponse(universe models.Universe) dtos.UniverseDtoResponse {
	return dtos.UniverseDtoResponse{
		Id:        universe.Id,
		Name:      universe.Name,
		CreatedAt: universe.CreatedAt,
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
