package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbPlanet struct {
	Id        uuid.UUID
	Player    uuid.UUID
	Name      string
	Homeworld bool

	Galaxy      int
	SolarSystem int
	Position    int

	Fields int

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int

	BuildingAction *uuid.UUID
}

func (p DbPlanet) ToDomain() models.Planet {
	return models.Planet{
		Id:        p.Id,
		Player:    p.Player,
		Name:      p.Name,
		Homeworld: p.Homeworld,
		Coordinate: models.Coordinate{
			Galaxy:      p.Galaxy,
			SolarSystem: p.SolarSystem,
			Position:    p.Position,
		},
		Fields:    p.Fields,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Version:   p.Version,
	}
}
