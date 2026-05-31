package models

import (
	"github.com/google/uuid"
)

type BuildingResourceStorage struct {
	Resource uuid.UUID
	Base     int
	Scale    float64
	Progress float64
}
