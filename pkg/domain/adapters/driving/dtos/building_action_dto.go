package dtos

import (
	"time"

	"github.com/google/uuid"
)

type BuildingActionDtoRequest struct {
	Building uuid.UUID `json:"building" format:"uuid"`
}

// TODO: Should also include the effects
type BuildingActionDtoResponse struct {
	Id           uuid.UUID `json:"id" format:"uuid"`
	Planet       uuid.UUID `json:"planet" format:"uuid"`
	Building     uuid.UUID `json:"building" format:"uuid"`
	CurrentLevel int       `json:"currentLevel"`
	DesiredLevel int       `json:"desiredLevel"`
	CreatedAt    time.Time `json:"createdAt" format:"date-time"`
	CompletedAt  time.Time `json:"completedAt" format:"date-time"`
}
