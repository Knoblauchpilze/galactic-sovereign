package dtos

import (
	"time"

	"github.com/google/uuid"
)

type BuildingActionDtoRequest struct {
	Building uuid.UUID `json:"building" format:"uuid"`
}

type BuildingActionDtoResponse struct {
	Id           uuid.UUID `json:"id" format:"uuid"`
	Planet       uuid.UUID `json:"planet" format:"uuid"`
	Building     uuid.UUID `json:"building" format:"uuid"`
	CurrentLevel int       `json:"currentLevel"`
	DesiredLevel int       `json:"desiredLevel"`

	CreatedAt   time.Time `json:"createdAt" format:"date-time"`
	CompletedAt time.Time `json:"completedAt" format:"date-time"`

	Costs       []BuildingActionCostDtoResponse       `json:"resources"`
	Storages    []BuildingActionStorageDtoResponse    `json:"storages"`
	Productions []BuildingActionProductionDtoResponse `json:"productions"`
}

type BuildingActionCostDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid"`
	Amount   int       `json:"amount"`
}

type BuildingActionStorageDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid"`
	Storage  int       `json:"storage"`
}

type BuildingActionProductionDtoResponse struct {
	Resource   uuid.UUID `json:"resource" format:"uuid"`
	Production int       `json:"production"`
}
