package models

import (
	"time"

	"github.com/google/uuid"
)

type BuildingAction struct {
	Id           uuid.UUID
	Planet       uuid.UUID
	Building     uuid.UUID
	CurrentLevel int
	DesiredLevel int

	CreatedAt   time.Time
	CompletedAt time.Time

	Version int

	Costs       []BuildingActionCost
	Storages    []BuildingActionResourceStorage
	Productions []BuildingActionResourceProduction
}

type BuildingActionCost struct {
	Resource uuid.UUID
	Amount   int
}

type BuildingActionResourceStorage struct {
	Resource uuid.UUID
	Storage  int
}

type BuildingActionResourceProduction struct {
	Resource   uuid.UUID
	Production int
}
