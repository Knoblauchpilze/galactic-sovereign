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
	}
}
