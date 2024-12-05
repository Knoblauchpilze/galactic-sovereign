package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	Id   uuid.UUID
	Name string

	StartAmount     int
	StartProduction int
	StartStorage    int

	CreatedAt time.Time
	UpdatedAt time.Time
}
