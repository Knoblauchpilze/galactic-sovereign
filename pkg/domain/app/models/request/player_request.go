package request

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type PlayerCreationRequest struct {
	ApiUser  uuid.UUID
	Universe uuid.UUID
	Name     string
}

func FromPlayerCreationRequest(player PlayerCreationRequest) models.Player {
	t := time.Now()
	return models.Player{
		Id:       uuid.New(),
		ApiUser:  player.ApiUser,
		Universe: player.Universe,
		Name:     player.Name,

		CreatedAt: t,

		Version: 0,
	}
}
