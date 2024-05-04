package persistence

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	Id      uuid.UUID
	Key     uuid.UUID
	ApiUser uuid.UUID

	ValidUntil time.Time
}
