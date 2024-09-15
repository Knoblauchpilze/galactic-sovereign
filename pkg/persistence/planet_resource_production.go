package persistence

import (
	"time"

	"github.com/google/uuid"
)

type PlanetResourceProduction struct {
	Planet     uuid.UUID
	Resource   uuid.UUID
	Building   *uuid.UUID
	Production int

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int
}
