package models

import (
	"time"

	"github.com/google/uuid"
)

type PlanetResourceProduction struct {
	Resource   uuid.UUID
	Building   *uuid.UUID
	Production int

	CreatedAt time.Time
	UpdatedAt time.Time
}
