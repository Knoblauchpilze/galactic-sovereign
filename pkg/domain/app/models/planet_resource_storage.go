package models

import (
	"time"

	"github.com/google/uuid"
)

type PlanetResourceStorage struct {
	Resource uuid.UUID
	Storage  int

	CreatedAt time.Time
	UpdatedAt time.Time
}
