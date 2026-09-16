package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbSolarSystem struct {
	Universe uuid.UUID
	Galaxy   int
	Number   int
	Orbits   int
}

func (ss DbSolarSystem) ToDomain() models.SolarSystem {
	return models.SolarSystem{
		Universe: ss.Universe,
		Galaxy:   ss.Galaxy,
		Number:   ss.Number,
		Orbits:   ss.Orbits,
	}
}
