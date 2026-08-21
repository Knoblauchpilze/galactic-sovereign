package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbPlayerPlanet struct {
	Id uuid.UUID

	Name string

	Galaxy      int
	SolarSystem int
	Position    int
}

func (p DbPlayerPlanet) ToDomain() models.PlayerPlanet {
	return models.PlayerPlanet{
		Id:   p.Id,
		Name: p.Name,
		Coordinate: models.Coordinate{
			Galaxy:      p.Galaxy,
			SolarSystem: p.SolarSystem,
			Position:    p.Position,
		},
	}
}
