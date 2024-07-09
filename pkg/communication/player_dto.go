package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type PlayerDtoRequest struct {
	ApiUser  uuid.UUID `json:"api_user"`
	Universe uuid.UUID `json:"universe"`
	Name     string    `json:"name" form:"name"`
}

type PlayerDtoResponse struct {
	Id       uuid.UUID `json:"id"`
	ApiUser  uuid.UUID `json:"api_user"`
	Universe uuid.UUID `json:"universe"`
	Name     string    `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
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
