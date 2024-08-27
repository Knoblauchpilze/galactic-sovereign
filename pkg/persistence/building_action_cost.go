package persistence

import (
	"github.com/google/uuid"
)

type BuildingActionCost struct {
	Action   uuid.UUID
	Resource uuid.UUID
	Amount   int
}
