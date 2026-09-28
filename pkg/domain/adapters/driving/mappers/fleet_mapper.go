package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
)

const (
	invalidFleetMission = models.FleetMission("")
)

func ToFleetCreationRequest(
	planetId uuid.UUID,
	dto dtos.FleetDtoRequest,
) (request.FleetCreationRequest, error) {
	mission, err := toFleetMission(dto.Mission)
	if err != nil {
		return request.FleetCreationRequest{}, err
	}

	req := request.FleetCreationRequest{
		Planet:  planetId,
		Mission: mission,
		Destination: request.FleetDestinationRequest{
			Galaxy:      dto.Destination.Galaxy,
			SolarSystem: dto.Destination.SolarSystem,
			Position:    dto.Destination.Position,
		},
		Ships: toFleetShipsRequest(dto.Ships),
	}

	return req, nil
}

func toFleetMission(mission dtos.FleetMissionDto) (models.FleetMission, error) {
	switch mission {
	case dtos.MissionColonize:
		return models.MissionColonize, nil
	default:
		return invalidFleetMission, ErrInvalidEnumMapping
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
