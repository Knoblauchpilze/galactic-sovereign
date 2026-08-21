package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	homeworldDefaultName string = "homeworld"
	planetDefaultName    string = "colony"
)

type PlayerPlanet struct {
	Id         uuid.UUID
	Name       string
	Coordinate Coordinate
}

type Player struct {
	Id       uuid.UUID
	ApiUser  uuid.UUID
	Universe uuid.UUID
	Name     string

	CreatedAt time.Time

	Version int

	Homeworld uuid.UUID
	Planets   []PlayerPlanet
}

func (p *Player) CreateHomeworld(
	universe Universe,
) Planet {
	planet := universe.CreatePlanet(p.Id, true)

	p.Homeworld = planet.Id
	p.Planets = []PlayerPlanet{
		{
			Id:         planet.Id,
			Name:       planet.Name,
			Coordinate: planet.Coordinate,
		},
	}

	return planet
}

func (p *Player) Colonize(
	universe Universe,
) Planet {
	planet := universe.CreatePlanet(p.Id, false)

	p.Planets = append(p.Planets, PlayerPlanet{
		Id:         planet.Id,
		Name:       planet.Name,
		Coordinate: planet.Coordinate,
	})

	return planet
}
