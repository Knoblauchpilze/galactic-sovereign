package models

import (
	"github.com/google/uuid"
)

type BuildingActionResourceProduction struct {
	Resource   uuid.UUID
	Production int
}
