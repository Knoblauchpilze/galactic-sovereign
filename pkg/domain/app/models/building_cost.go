package models

import (
	"github.com/google/uuid"
)

type BuildingCost struct {
	Resource uuid.UUID
	Cost     int
	Progress float64
}
