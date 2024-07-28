package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	Id   uuid.UUID
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
}
