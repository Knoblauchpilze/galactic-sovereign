package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Building struct {
	Id   uuid.UUID
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
}
