package models

import (
	"time"

	"github.com/google/uuid"
)

type Universe struct {
	Id   uuid.UUID
	Name string

	CreatedAt time.Time

	Version int

	Resources []Resource
}
