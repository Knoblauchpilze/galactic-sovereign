package models

import (
	"time"

	"github.com/google/uuid"
)

type BuildingAction struct {
	Id       uuid.UUID
	Building uuid.UUID

	DesiredLevel int

	CreatedAt   time.Time
	CompletedAt time.Time

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
