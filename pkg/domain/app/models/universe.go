package models

import (
	"time"

	"github.com/google/uuid"
)

// TODO: This should include the buildings and their effects
type Universe struct {
	Id   uuid.UUID
	Name string

	CreatedAt time.Time

	Version int

	Buildings []Building
}
