package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	homeworldDefaultName string = "homeworld"
	planetDefaultName    string = "colony"
)

type Player struct {
	Id       uuid.UUID
	ApiUser  uuid.UUID
	Universe uuid.UUID
	Name     string

	CreatedAt time.Time

	Version int

	Homeworld uuid.UUID
	Planets   []uuid.UUID
}

func (p *Player) CreateHomeworld(
	universe Universe,
) Planet {
	planet := universe.CreatePlanet(p.Id, true)

	p.Homeworld = planet.Id
	p.Planets = []uuid.UUID{planet.Id}

	return planet
}

func (p *Player) Colonize(
	universe Universe,
) Planet {
	planet := universe.CreatePlanet(p.Id, false)

	p.Planets = append(p.Planets, planet.Id)

	return planet
}
