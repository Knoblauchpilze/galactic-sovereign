package persistence

import (
	"time"

	"github.com/google/uuid"
)

type BuildingAction struct {
	Id           uuid.UUID
	Planet       uuid.UUID
	Building     uuid.UUID
	CurrentLevel int
	DesiredLevel int
	CreatedAt    time.Time
	CompletedAt  time.Time
}
