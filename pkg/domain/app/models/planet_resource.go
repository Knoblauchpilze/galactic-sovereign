package models

import (
	"time"

	"github.com/google/uuid"
)

type PlanetResource struct {
	Resource uuid.UUID
	Amount   float64

	CreatedAt time.Time
	UpdatedAt time.Time
}
