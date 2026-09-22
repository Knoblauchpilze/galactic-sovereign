package dtos

import (
	"time"

	"github.com/google/uuid"
)

type FleetDtoRequest struct {
	Ships []FleetShipDtoRequest `json:"ships" binding:"required,min=1,dive"`
}

type FleetShipDtoRequest struct {
	Ship  uuid.UUID `json:"ship" format:"uuid" binding:"required"`
	Count int       `json:"count" binding:"required,min=1" minimum:"1"`
}

type FleetDtoResponse struct {
	Id uuid.UUID `json:"id" format:"uuid" binding:"required"`

	CreatedAt time.Time `json:"created_at" format:"date-time" binding:"required"`

	Ships []FleetShipDtoResponse `json:"costs" binding:"required"`
}

type FleetShipDtoResponse struct {
	Ship  uuid.UUID `json:"ship" format:"uuid" binding:"required"`
	Count int       `json:"count" binding:"required"`
}
