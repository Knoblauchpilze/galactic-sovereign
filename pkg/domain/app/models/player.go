package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	homeworldDefaultName string = "homeworld"
)

type Player struct {
	Id       uuid.UUID
	ApiUser  uuid.UUID
	Universe uuid.UUID
	Name     string

	CreatedAt time.Time

	Version int

	Planets []uuid.UUID
}

func (p *Player) CreateHomeworld() Planet {
	createdAt := time.Now()

	planet := Planet{
		Id:             uuid.New(),
		Player:         p.Id,
		Name:           homeworldDefaultName,
		Homeworld:      true,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
		Version:        0,
		Resources:      []PlanetResource{},
		Storages:       []PlanetResourceStorage{},
		Productions:    []PlanetResourceProduction{},
		Buildings:      []PlanetBuilding{},
		BuildingAction: nil,
	}

	p.Planets = []uuid.UUID{planet.Id}

	return planet
}
