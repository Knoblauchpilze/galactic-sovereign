package models

import (
	"time"

	"github.com/google/uuid"
)

type ShipAction struct {
	Id   uuid.UUID
	Ship uuid.UUID

	Count int

	CreatedAt        time.Time
	NextCompletionAt time.Time
	CompletedAt      time.Time

	Costs []ShipActionCost
}

type ShipActionCost struct {
	Resource uuid.UUID
	Amount   int
}
