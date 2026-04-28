package communication

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlayerDtoRequest struct {
	ApiUser  uuid.UUID `json:"api_user" format:"uuid"`
	Universe uuid.UUID `json:"universe" format:"uuid"`
	Name     string    `json:"name" form:"name"`
}

type PlayerDtoResponse struct {
	Id       uuid.UUID `json:"id" format:"uuid"`
	ApiUser  uuid.UUID `json:"api_user" format:"uuid"`
	Universe uuid.UUID `json:"universe" format:"uuid"`
	Name     string    `json:"name"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
}

func FromPlayerDtoRequest(player PlayerDtoRequest) persistence.Player {
	t := time.Now()
	return persistence.Player{
		Id:       uuid.New(),
		ApiUser:  player.ApiUser,
		Universe: player.Universe,
		Name:     player.Name,

		CreatedAt: t,
		UpdatedAt: t,
	}
}

func ToPlayerDtoResponse(player persistence.Player) PlayerDtoResponse {
	return PlayerDtoResponse{
		Id:       player.Id,
		ApiUser:  player.ApiUser,
		Universe: player.Universe,
		Name:     player.Name,

		CreatedAt: player.CreatedAt,
	}
}
