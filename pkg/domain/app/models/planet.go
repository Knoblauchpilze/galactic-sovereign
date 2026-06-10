package models

import (
	"time"

	"github.com/google/uuid"
)

type Planet struct {
	Id        uuid.UUID
	Player    uuid.UUID
	Name      string
	Homeworld bool

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int

	Resources   []PlanetResource
	Storages    []PlanetResourceStorage
	Productions []PlanetResourceProduction

	Buildings []PlanetBuilding

	BuildingAction *uuid.UUID
}

type PlanetResource struct {
	Resource uuid.UUID
	Amount   float64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PlanetResourceStorage struct {
	Resource uuid.UUID
	Storage  int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PlanetResourceProduction struct {
	Resource   uuid.UUID
	Building   *uuid.UUID
	Production int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PlanetBuilding struct {
	Building uuid.UUID
	Level    int

	CreatedAt time.Time
	UpdatedAt time.Time
}
