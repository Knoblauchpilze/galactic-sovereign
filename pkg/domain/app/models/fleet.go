package models

import (
	"time"

	"github.com/google/uuid"
)

type Fleet struct {
	Id          uuid.UUID
	Player      uuid.UUID
	Source      uuid.UUID
	Destination Coordinate
	Ships       []FleetShip
	CreatedAt   time.Time
	ArrivalAt   time.Time
	ReturnAt    time.Time
	UpdatedAt   time.Time
	Version     int
}

type FleetShip struct {
	Ship  uuid.UUID
	Count int
}
