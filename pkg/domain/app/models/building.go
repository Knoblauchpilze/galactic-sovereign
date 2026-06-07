package models

import (
	"time"

	"github.com/google/uuid"
)

// TODO: A new function should create a building action for a certain level
type Building struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time

	Costs       []BuildingCost
	Productions []BuildingResourceProduction
	Storages    []BuildingResourceStorage
}

type BuildingCost struct {
	Resource uuid.UUID
	Cost     int
	Progress float64
}

type BuildingResourceProduction struct {
	Resource uuid.UUID
	Base     int
	Progress float64
}

type BuildingResourceStorage struct {
	Resource uuid.UUID
	Base     int
	Scale    float64
	Progress float64
}
