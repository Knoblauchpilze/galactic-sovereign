package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbBuilding struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
}

func (b DbBuilding) ToDomain() models.Building {
	return models.Building{
		Id:        b.Id,
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
	}
}
