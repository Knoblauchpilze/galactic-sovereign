package models

import "github.com/google/uuid"

type SolarSystem struct {
	Universe uuid.UUID
	Galaxy   int
	Number   int
	Orbits   int
	Planets  []SolarSystemPlanet
}

type SolarSystemPlanet struct {
	Id        uuid.UUID
	Player    uuid.UUID
	Name      string
	Homeworld bool
	Position  int
}
