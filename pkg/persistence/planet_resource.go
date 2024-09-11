package persistence

import (
	"time"

	"github.com/google/uuid"
)

type PlanetResource struct {
	Planet     uuid.UUID
	Resource   uuid.UUID
	Amount     float64
	Production int

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int
}
