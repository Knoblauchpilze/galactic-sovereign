package models

import (
	"time"

	"github.com/google/uuid"
)

type Building struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time

	Costs       []BuildingCost
	Productions []BuildingResourceProduction
}
