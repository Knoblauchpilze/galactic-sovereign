package dtos

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type FleetDtoRequest struct {
	Mission     FleetMissionDto            `json:"mission" binding:"required"`
	Destination FleetDestinationDtoRequest `json:"destination" binding:"required"`
	Ships       []FleetShipDtoRequest      `json:"ships" binding:"required,min=1,dive"`
}

type FleetDestinationDtoRequest struct {
	Galaxy      *int `json:"galaxy" binding:"required,gte=0" minimum:"0"`
	SolarSystem *int `json:"solar_system" binding:"required,gte=0" minimum:"0"`
	Position    *int `json:"position" binding:"required,gte=0" minimum:"0"`
}

type FleetShipDtoRequest struct {
	Ship  uuid.UUID `json:"ship" format:"uuid" binding:"required"`
	Count int       `json:"count" binding:"required,min=1" minimum:"1"`
}

type FleetDtoResponse struct {
	Id          uuid.UUID                   `json:"id" format:"uuid" binding:"required"`
	Player      uuid.UUID                   `json:"player" format:"uuid" binding:"required"`
	Source      uuid.UUID                   `json:"source" format:"uuid" binding:"required"`
	Mission     models.FleetMission         `json:"mission" binding:"required"`
	Destination FleetDestinationDtoResponse `json:"destination" binding:"required"`
	Schedule    FleetScheduleDtoResponse    `json:"schedule" binding:"required"`

	CreatedAt time.Time `json:"created_at" format:"date-time" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time" binding:"required"`

	Ships []FleetShipDtoResponse `json:"ships" binding:"required,min=1,dive"`
}

type FleetDestinationDtoResponse struct {
	Galaxy      int        `json:"galaxy" binding:"required" minimum:"0"`
	SolarSystem int        `json:"solar_system" binding:"required" minimum:"0"`
	Position    int        `json:"position" binding:"required" minimum:"0"`
	Target      *uuid.UUID `json:"target,omitempty" format:"uuid"`
}

type FleetScheduleDtoResponse struct {
	ArrivalAt time.Time          `json:"arrival_at" format:"date-time" binding:"required"`
	ReturnAt  time.Time          `json:"return_at" format:"date-time" binding:"required"`
	Status    models.FleetStatus `json:"status" binding:"required"`
}

type FleetShipDtoResponse struct {
	Ship  uuid.UUID `json:"ship" format:"uuid" binding:"required"`
	Count int       `json:"count" binding:"required"`
}
