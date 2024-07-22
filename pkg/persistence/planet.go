package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Planet struct {
	Id        uuid.UUID
	Player    uuid.UUID
	Name      string
	Homeworld bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
