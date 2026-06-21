package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbPlayer struct {
	Id       uuid.UUID
	ApiUser  uuid.UUID
	Universe uuid.UUID
	Name     string

	CreatedAt time.Time

	Version int

	Homeworld uuid.UUID
}

func (p DbPlayer) ToDomain() models.Player {
	return models.Player{
		Id:        p.Id,
		ApiUser:   p.ApiUser,
		Universe:  p.Universe,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		Version:   p.Version,
		Homeworld: p.Homeworld,
	}
}
