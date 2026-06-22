package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbUniverse struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	Version   int
}

func (p DbUniverse) ToDomain() models.Universe {
	return models.Universe{
		Id:        p.Id,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		Version:   p.Version,
	}
}
