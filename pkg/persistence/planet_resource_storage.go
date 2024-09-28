package persistence

import (
	"time"

	"github.com/google/uuid"
)

type PlanetResourceStorage struct {
	Planet   uuid.UUID
	Resource uuid.UUID
	Storage  int

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int
}
