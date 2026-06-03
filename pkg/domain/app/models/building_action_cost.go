package models

import (
	"github.com/google/uuid"
)

type BuildingActionCost struct {
	Resource uuid.UUID
	Amount   int
}
