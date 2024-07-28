package persistence

import (
	"time"

	"github.com/google/uuid"
)

type PlanetResource struct {
	Planet   uuid.UUID
	Resource uuid.UUID
	Amount   float64

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int
}
