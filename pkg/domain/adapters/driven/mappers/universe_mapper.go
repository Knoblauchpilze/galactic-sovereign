package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbUniverse struct {
	Id   uuid.UUID
	Name string

	CreatedAt time.Time

	Version int
}

func (u DbUniverse) ToDomain() models.Universe {
	return models.Universe{
		Id:        u.Id,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		Version:   u.Version,
	}
}
