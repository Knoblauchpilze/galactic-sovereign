package persistence

import (
	"time"

	"github.com/google/uuid"
)

type UserLimit struct {
	Id   uuid.UUID
	Name string
	User uuid.UUID

	Limits []Limit

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Limit struct {
	Id uuid.UUID

	Name  string
	Value string

	CreatedAt time.Time
	UpdatedAt time.Time
}
