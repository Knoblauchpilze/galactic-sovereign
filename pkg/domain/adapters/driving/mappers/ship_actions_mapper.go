package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
)

func ToShipActionCreationRequest(
	planetId uuid.UUID,
	dto dtos.ShipActionDtoRequest,
) request.ShipActionCreationRequest {
	return request.ShipActionCreationRequest{
		Planet: planetId,
		Ship:   dto.Ship,
		Count:  dto.Count,
	}
}

func ToShipActionResponse(action models.ShipAction) dtos.ShipActionDtoResponse {
	return dtos.ShipActionDtoResponse{
		Id:                 action.Id,
		Ship:               action.Ship,
		Count:              action.Count,
		CreatedAt:          action.CreatedAt,
		NextCompletionAt:   action.NextCompletionAt,
		UnitCompletionTime: toIso8601Duration(action.UnitCompletionTime),
		Costs:              toShipActionCostsResponse(action.Costs),
	}
}

func toShipActionCostResponse(
	cost models.ShipActionCost,
) dtos.ShipActionCostDtoResponse {
	return dtos.ShipActionCostDtoResponse{
		Resource: cost.Resource,
		Amount:   cost.Amount,
	}
}

func toShipActionCostsResponse(
	costs []models.ShipActionCost,
) []dtos.ShipActionCostDtoResponse {
	out := make([]dtos.ShipActionCostDtoResponse, 0, len(costs))

	for _, c := range costs {
		dto := toShipActionCostResponse(c)
		out = append(out, dto)
	}

	return out
}
