package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

const (
	invalidFleetMission = models.FleetMission("")
)

func ToFleetOrder(
	dto dtos.FleetDtoRequest,
) (models.FleetOrder, error) {
	mission, err := toFleetMission(dto.Mission)
	if err != nil {
		return models.FleetOrder{}, err
	}

	req := models.FleetOrder{
		Mission: mission,
		Destination: models.Coordinate{
			// Should always be not nil as the API specs requires it
			// and it's enforced by the validator (in the controller)
			Galaxy:      *dto.Destination.Galaxy,
			SolarSystem: *dto.Destination.SolarSystem,
			Position:    *dto.Destination.Position,
		},
		Ships: toFleetShips(dto.Ships),
	}

	return req, nil
}

func toFleetMission(mission dtos.FleetMissionDto) (models.FleetMission, error) {
	switch mission {
	case dtos.MissionColonize:
		return models.MissionColonize, nil
	case dtos.MissionTransport:
		return models.MissionTransport, nil
	default:
		return invalidFleetMission, ErrInvalidEnumMapping
	}
}

func toFleetShip(
	ship dtos.FleetShipDtoRequest,
) models.FleetShip {
	return models.FleetShip{
		Ship:  ship.Ship,
		Count: ship.Count,
	}
}

func toFleetShips(
	ships []dtos.FleetShipDtoRequest,
) []models.FleetShip {
	out := make([]models.FleetShip, 0, len(ships))

	for _, s := range ships {
		dto := toFleetShip(s)
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
