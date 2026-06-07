package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
)

func ToBuildingActionCreationRequest(
	planetId uuid.UUID,
	dto dtos.BuildingActionDtoRequest,
) request.BuildingActionCreationRequest {
	return request.BuildingActionCreationRequest{
		Planet:   planetId,
		Building: dto.Building,
	}
}

func ToBuildingActionResponse(action models.BuildingAction) dtos.BuildingActionDtoResponse {
	return dtos.BuildingActionDtoResponse{
		Id:           action.Id,
		Planet:       action.Planet,
		Building:     action.Building,
		CurrentLevel: action.CurrentLevel,
		DesiredLevel: action.DesiredLevel,
		CreatedAt:    action.CreatedAt,
		CompletedAt:  action.CompletedAt,
		Costs:        toBuildingActionCostsResponse(action.Costs),
		Storages:     toBuildingActionStoragesResponse(action.Storages),
		Productions:  toBuildingActionProductionsResponse(action.Productions),
	}
}

func toBuildingActionCostResponse(
	cost models.BuildingActionCost,
) dtos.BuildingActionCostDtoResponse {
	return dtos.BuildingActionCostDtoResponse{
		Resource: cost.Resource,
		Amount:   cost.Amount,
	}
}

func toBuildingActionCostsResponse(
	costs []models.BuildingActionCost,
) []dtos.BuildingActionCostDtoResponse {
	out := make([]dtos.BuildingActionCostDtoResponse, 0, len(costs))

	for _, c := range costs {
		dto := toBuildingActionCostResponse(c)
		out = append(out, dto)
	}

	return out
}

func toBuildingActionStorageResponse(
	storage models.BuildingActionResourceStorage,
) dtos.BuildingActionStorageDtoResponse {
	return dtos.BuildingActionStorageDtoResponse{
		Resource: storage.Resource,
		Storage:  storage.Storage,
	}
}

func toBuildingActionStoragesResponse(
	storages []models.BuildingActionResourceStorage,
) []dtos.BuildingActionStorageDtoResponse {
	out := make([]dtos.BuildingActionStorageDtoResponse, 0, len(storages))

	for _, s := range storages {
		dto := toBuildingActionStorageResponse(s)
		out = append(out, dto)
	}

	return out
}

func toBuildingActionProductionResponse(
	production models.BuildingActionResourceProduction,
) dtos.BuildingActionProductionDtoResponse {
	return dtos.BuildingActionProductionDtoResponse{
		Resource:   production.Resource,
		Production: production.Production,
	}
}

func toBuildingActionProductionsResponse(
	productions []models.BuildingActionResourceProduction,
) []dtos.BuildingActionProductionDtoResponse {
	out := make([]dtos.BuildingActionProductionDtoResponse, 0, len(productions))

	for _, p := range productions {
		dto := toBuildingActionProductionResponse(p)
		out = append(out, dto)
	}

	return out
}
