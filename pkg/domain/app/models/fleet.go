package models

import (
	"time"

	"github.com/google/uuid"
)

type Fleet struct {
	Id          uuid.UUID
	Player      uuid.UUID
	Source      uuid.UUID
	Destination FleetDestination
	Ships       []FleetShip
	CreatedAt   time.Time
	ArrivalAt   time.Time
	ReturnAt    time.Time
	UpdatedAt   time.Time
	Version     int
}

type FleetDestination struct {
	Coordinate Coordinate
	Target     *uuid.UUID
}

type FleetShip struct {
	Ship  uuid.UUID
	Count int
}
