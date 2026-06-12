package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
)

func ToPlayerCreationRequest(dto dtos.PlayerDtoRequest) request.PlayerCreationRequest {
	return request.PlayerCreationRequest{
		ApiUser:  dto.ApiUser,
		Universe: dto.Universe,
		Name:     dto.Name,
	}
}

func ToPlayerResponse(player models.Player) dtos.PlayerDtoResponse {
	return dtos.PlayerDtoResponse{
		Id:        player.Id,
		ApiUser:   player.ApiUser,
		Universe:  player.Universe,
		Name:      player.Name,
		CreatedAt: player.CreatedAt,
		Planets:   player.Planets,
	}
}

func ToPlayersResponse(players []models.Player) []dtos.PlayerDtoResponse {
	out := make([]dtos.PlayerDtoResponse, 0, len(players))

	for _, p := range players {
		dto := ToPlayerResponse(p)
		out = append(out, dto)
	}

	return out
}
