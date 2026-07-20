package models

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	Id   uuid.UUID
	Name string

	StartAmount     int
	StartProduction int
	StartStorage    int

	// Defines how long each unit of a resource takes to be transformed into
	// a building. A value of 1 means that if a building costs 3 unit of the
	// resource it will take 3 hours to be built.
	BuildTimeHoursPerUnit float64

	CreatedAt time.Time
}
