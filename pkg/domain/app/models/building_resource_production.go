package models

import (
	"github.com/google/uuid"
)

type BuildingResourceProduction struct {
	Resource uuid.UUID
	Base     int
	Progress float64
}
