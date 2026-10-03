package models

import "github.com/google/uuid"

type FleetOrigin struct {
	Source     uuid.UUID
	Coordinate Coordinate
	Player     uuid.UUID
}
