package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

func ToSolarSystemResponse(solarSystem models.SolarSystem) dtos.SolarSystemDtoResponse {
	return dtos.SolarSystemDtoResponse{
		Universe: solarSystem.Universe,
		Galaxy:   solarSystem.Galaxy,
		Number:   solarSystem.Number,
		Orbits:   solarSystem.Orbits,
		Planets:  toSolarSystemPlanetsResponse(solarSystem.Planets),
	}
}

func toSolarSystemPlanetResponse(
	planet models.SolarSystemPlanet,
) dtos.SolarSystemPlanetDtoResponse {
	return dtos.SolarSystemPlanetDtoResponse{
		Id:        planet.Id,
		Player:    planet.Player,
		Name:      planet.Name,
		Homeworld: planet.Homeworld,
		Position:  planet.Position,
	}
}

func toSolarSystemPlanetsResponse(
	planets []models.SolarSystemPlanet,
) []dtos.SolarSystemPlanetDtoResponse {
	out := make([]dtos.SolarSystemPlanetDtoResponse, 0, len(planets))

	for _, p := range planets {
		dto := toSolarSystemPlanetResponse(p)
		out = append(out, dto)
	}

	return out
}
