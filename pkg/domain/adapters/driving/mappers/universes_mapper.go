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
		Resources: toResourcesResponse(universe.Resources),
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

func toResourceResponse(
	resource models.Resource,
) dtos.ResourceDtoResponse {
	return dtos.ResourceDtoResponse{
		Id:              resource.Id,
		Name:            resource.Name,
		StartAmount:     resource.StartAmount,
		StartProduction: resource.StartProduction,
		StartStorage:    resource.StartStorage,
		CreatedAt:       resource.CreatedAt,
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
