package persistence

import (
	"github.com/google/uuid"
)

type BuildingResourceProduction struct {
	Building uuid.UUID
	Resource uuid.UUID
	Base     int
	Progress float64
}
