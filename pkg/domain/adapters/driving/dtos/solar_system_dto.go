package dtos

import (
	"github.com/google/uuid"
)

type SolarSystemDtoResponse struct {
	Universe uuid.UUID `json:"id" format:"uuid" binding:"required"`
	Galaxy   int       `json:"galaxy" binding:"required" minimum:"0"`
	Number   int       `json:"number" binding:"required" minimum:"0"`
	Orbits   int       `json:"position" binding:"required" minimum:"0"`
	Planets  []SolarSystemPlanetDtoResponse
}

type SolarSystemPlanetDtoResponse struct {
	Id        uuid.UUID `json:"id" format:"uuid" binding:"required"`
	Player    uuid.UUID `json:"player" format:"uuid" binding:"required"`
	Name      string    `json:"name" binding:"required" example:"colony"`
	Homeworld bool      `json:"homeworld" binding:"required"`
	Position  int       `json:"position" binding:"required" minimum:"0"`
}
