package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbShipAction struct {
	Id    uuid.UUID
	Ship  uuid.UUID
	Count int

	CreatedAt        time.Time
	NextCompletionAt time.Time
	CompletedAt      time.Time
}

func (a DbShipAction) ToDomain() models.ShipAction {
	return models.ShipAction{
		Id:    a.Id,
		Ship:  a.Ship,
		Count: a.Count,

		CreatedAt:        a.CreatedAt,
		NextCompletionAt: a.NextCompletionAt,
		CompletedAt:      a.CompletedAt,
	}
}
