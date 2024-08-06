package persistence

import (
	"time"

	"github.com/google/uuid"
)

type PlanetBuilding struct {
	Planet   uuid.UUID
	Building uuid.UUID
	Level    int

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int
}
