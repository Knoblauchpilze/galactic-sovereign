package dtos

import (
	"time"

	"github.com/google/uuid"
)

type BuildingActionDtoRequest struct {
	Building uuid.UUID `json:"building" format:"uuid" binding:"required"`
}

type BuildingActionDtoResponse struct {
	Id           uuid.UUID `json:"id" format:"uuid" binding:"required"`
	Planet       uuid.UUID `json:"planet" format:"uuid" binding:"required"`
	Building     uuid.UUID `json:"building" format:"uuid" binding:"required"`
	DesiredLevel int       `json:"desired_level" binding:"required"`

	CreatedAt   time.Time `json:"created_at" format:"date-time" binding:"required"`
	CompletedAt time.Time `json:"completed_at" format:"date-time" binding:"required"`

	Costs       []BuildingActionCostDtoResponse       `json:"resources" binding:"required"`
	Storages    []BuildingActionStorageDtoResponse    `json:"storages" binding:"required"`
	Productions []BuildingActionProductionDtoResponse `json:"productions" binding:"required"`
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
