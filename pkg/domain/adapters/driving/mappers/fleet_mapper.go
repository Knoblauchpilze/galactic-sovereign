package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
)

func ToFleetCreationRequest(
	planetId uuid.UUID,
	dto dtos.FleetDtoRequest,
) request.FleetCreationRequest {
	return request.FleetCreationRequest{
		Planet: planetId,
		Ships:  toFleetShipsRequest(dto.Ships),
	}
}

func toFleetShipRequest(
	ship dtos.FleetShipDtoRequest,
) request.FleetShipRequest {
	return request.FleetShipRequest{
		Ship:  ship.Ship,
		Count: ship.Count,
	}
}

func toFleetShipsRequest(
	ships []dtos.FleetShipDtoRequest,
) []request.FleetShipRequest {
	out := make([]request.FleetShipRequest, 0, len(ships))

	for _, s := range ships {
		dto := toFleetShipRequest(s)
		out = append(out, dto)
	}

	return out
}

func ToFleetResponse(fleet models.Fleet) dtos.FleetDtoResponse {
	return dtos.FleetDtoResponse{
		Id:        fleet.Id,
		CreatedAt: fleet.CreatedAt,
		Ships:     toFleetShipsResponse(fleet.Ships),
	}
}

func toFleetShipResponse(
	ship models.FleetShip,
) dtos.FleetShipDtoResponse {
	return dtos.FleetShipDtoResponse{
		Ship:  ship.Ship,
		Count: ship.Count,
	}
}

func toFleetShipsResponse(
	ships []models.FleetShip,
) []dtos.FleetShipDtoResponse {
	out := make([]dtos.FleetShipDtoResponse, 0, len(ships))

	for _, s := range ships {
		dto := toFleetShipResponse(s)
		out = append(out, dto)
	}

	return out
}
