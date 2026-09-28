package dtos

import (
	"time"

	"github.com/google/uuid"
)

type FleetDtoRequest struct {
	Mission     FleetMissionDto            `json:"mission" binding:"required"`
	Destination FleetDestinationDtoRequest `json:"destination" binding:"required"`
	Ships       []FleetShipDtoRequest      `json:"ships" binding:"required,min=1,dive"`
}

type FleetDestinationDtoRequest struct {
	Galaxy      int `json:"galaxy" binding:"required" minimum:"0"`
	SolarSystem int `json:"solar_system" binding:"required" minimum:"0"`
	Position    int `json:"position" binding:"required" minimum:"0"`
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
