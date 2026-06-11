package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PlanetDtoRequest struct {
	Player uuid.UUID `json:"player" format:"uuid"`
	Name   string    `json:"name"`
}

type PlanetDtoResponse struct {
	Id        uuid.UUID `json:"id" format:"uuid"`
	Player    uuid.UUID `json:"player" format:"uuid"`
	Name      string    `json:"name"`
	Homeworld bool      `json:"homeworld"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`

	Resources   []PlanetResourceDtoResponse           `json:"resources"`
	Storages    []PlanetResourceStorageDtoResponse    `json:"storages"`
	Productions []PlanetResourceProductionDtoResponse `json:"productions"`
	Buildings   []PlanetBuildingDtoResponse           `json:"buildings"`

	BuildingAction *uuid.UUID `json:"building_action,omitempty"`
}

type PlanetResourceDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid"`
	Amount   float64   `json:"amount"`
}

type PlanetResourceStorageDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid"`
	Storage  int       `json:"storage"`
}

type PlanetResourceProductionDtoResponse struct {
	Building   *uuid.UUID `json:"building,omitempty" format:"uuid"`
	Resource   uuid.UUID  `json:"resource" format:"uuid"`
	Production int        `json:"production"`
}

type PlanetBuildingDtoResponse struct {
	Building  uuid.UUID `json:"building" format:"uuid"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`
}
