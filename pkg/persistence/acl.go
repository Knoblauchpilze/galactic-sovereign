package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Acl struct {
	Id   uuid.UUID
	User uuid.UUID

	Resource    string
	Permissions []string

	CreatedAt time.Time
	UpdatedAt time.Time
}
