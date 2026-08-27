package models

import (
	"time"

	"github.com/google/uuid"
)

type Ship struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time

	Costs []ShipCost
}

type ShipCost struct {
	Resource uuid.UUID
	Cost     int
}
