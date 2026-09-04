package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
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
