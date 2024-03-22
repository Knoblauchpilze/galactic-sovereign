package persistence

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Email    string
	Password string

	CreatedAt  time.Time
	Updated_at time.Time
}
