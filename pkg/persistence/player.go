package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	Id       uuid.UUID
	ApiUser  uuid.UUID
	Universe uuid.UUID
	Name     string

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int
}
