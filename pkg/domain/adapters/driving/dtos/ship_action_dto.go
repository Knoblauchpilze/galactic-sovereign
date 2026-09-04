package dtos

import (
	"time"

	"github.com/google/uuid"
)

type ShipActionDtoRequest struct {
	Ship  uuid.UUID `json:"ship" format:"uuid" binding:"required"`
	Count int       `json:"count" binding:"required,min=1" minimum:"1"`
}

type ShipActionDtoResponse struct {
	Id    uuid.UUID `json:"id" format:"uuid" binding:"required"`
	Ship  uuid.UUID `json:"ship" format:"uuid" binding:"required"`
	Count int       `json:"count" binding:"required"`

	CreatedAt          time.Time     `json:"created_at" format:"date-time" binding:"required"`
	NextCompletionAt   time.Time     `json:"next_completion_at" format:"date-time" binding:"required"`
	UnitCompletionTime time.Duration `json:"unit_completion_time" format:"duration" binding:"required"`

	Costs []ShipActionCostDtoResponse `json:"costs" binding:"required"`
}

type ShipActionCostDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid"`
	Amount   int       `json:"amount"`
}
