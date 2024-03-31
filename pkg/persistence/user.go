package persistence

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Email    string
	Password string

	ApiKeys []uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time

	// https://stackoverflow.com/questions/129329/optimistic-vs-pessimistic-locking
	Version int
}
