package persistence

import (
	"github.com/google/uuid"
)

type BuildingResourceStorage struct {
	Building uuid.UUID
	Resource uuid.UUID
	Base     int
	Scale    float64
	Progress float64
}
