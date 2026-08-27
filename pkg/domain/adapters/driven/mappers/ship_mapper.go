package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbShip struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
}

func (s DbShip) ToDomain() models.Ship {
	return models.Ship{
		Id:        s.Id,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
	}
}
