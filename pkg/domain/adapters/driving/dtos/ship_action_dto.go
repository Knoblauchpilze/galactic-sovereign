package dtos

import (
	"github.com/google/uuid"
)

type ShipActionDtoRequest struct {
	Ship  uuid.UUID `json:"ship" format:"uuid" binding:"required"`
	Count int       `json:"count" binding:"required,min=1" minimum:"1"`
}
