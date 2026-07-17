package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PlanetDtoResponse struct {
	Id         uuid.UUID             `json:"id" format:"uuid" binding:"required"`
	Player     uuid.UUID             `json:"player" format:"uuid" binding:"required"`
	Name       string                `json:"name" example:"colony" binding:"required"`
	Homeworld  bool                  `json:"homeworld" binding:"required"`
	Coordinate CoordinateDtoResponse `json:"coordinate" binding:"required"`
	Fields     int                   `json:"fields" binding:"required" minimum:"1"`

	CreatedAt time.Time `json:"created_at" format:"date-time" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time" binding:"required"`

	Resources   []PlanetResourceDtoResponse           `json:"resources" binding:"required"`
	Storages    []PlanetResourceStorageDtoResponse    `json:"storages" binding:"required"`
	Productions []PlanetResourceProductionDtoResponse `json:"productions" binding:"required"`
	Buildings   []PlanetBuildingDtoResponse           `json:"buildings" binding:"required"`

	BuildingAction *BuildingActionDtoResponse `json:"building_action,omitempty"`
}

type CoordinateDtoResponse struct {
	Galaxy      int `json:"galaxy" binding:"required" minimum:"1"`
	SolarSystem int `json:"solar_system" binding:"required" minimum:"1"`
	Position    int `json:"position" binding:"required" minimum:"1"`
}

type PlanetResourceDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid" binding:"required"`
	Amount   float64   `json:"amount" binding:"required"`
}

type PlanetResourceStorageDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid" binding:"required"`
	Storage  int       `json:"storage" binding:"required"`
}

type PlanetResourceProductionDtoResponse struct {
	Building   *uuid.UUID `json:"building,omitempty" format:"uuid" binding:"required"`
	Resource   uuid.UUID  `json:"resource" format:"uuid" binding:"required"`
	Production int        `json:"production" binding:"required"`
}

type PlanetBuildingDtoResponse struct {
	Building uuid.UUID `json:"building" format:"uuid" binding:"required"`
	Level    int       `json:"level" binding:"required"`
}
