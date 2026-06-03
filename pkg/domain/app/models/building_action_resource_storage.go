package models

import (
	"github.com/google/uuid"
)

type BuildingActionResourceStorage struct {
	Resource uuid.UUID
	Storage  int
}
