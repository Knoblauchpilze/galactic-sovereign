package persistence

import (
	"github.com/google/uuid"
)

type BuildingActionResourceProduction struct {
	Action     uuid.UUID
	Resource   uuid.UUID
	Production int
}
