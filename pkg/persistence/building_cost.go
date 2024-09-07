package persistence

import (
	"github.com/google/uuid"
)

type BuildingCost struct {
	Building uuid.UUID
	Resource uuid.UUID
	Cost     int
	Progress float64
}
