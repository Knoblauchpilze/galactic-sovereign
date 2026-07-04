package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UniverseDtoRequest struct {
	Name string `json:"name" example:"aquarius" binding:"required"`
}

type UniverseDtoResponse struct {
	Id   uuid.UUID `json:"id" format:"uuid" binding:"required"`
	Name string    `json:"name" example:"oberon" binding:"required"`

	CreatedAt time.Time `json:"created_at" format:"date-time" binding:"required"`

	Resources []ResourceDtoResponse `json:"resources" binding:"required"`
	Buildings []BuildingDtoResponse `json:"buildings" binding:"required"`
}

type ResourceDtoResponse struct {
	Id   uuid.UUID `json:"id" format:"uuid" binding:"required"`
	Name string    `json:"name" example:"metal" binding:"required"`

	StartAmount     int `json:"start_amount" binding:"required"`
	StartProduction int `json:"start_production" binding:"required"`
	StartStorage    int `json:"start_storage" binding:"required"`

	CreatedAt time.Time `json:"created_at" format:"date-time" binding:"required"`
}

type BuildingDtoResponse struct {
	Id        uuid.UUID `json:"id" format:"uuid" binding:"required"`
	Name      string    `json:"name" example:"metal mine" binding:"required"`
	CreatedAt time.Time `json:"created_at" format:"date-time" binding:"required"`

	Costs       []BuildingCostDtoResponse               `json:"costs" binding:"required"`
	Productions []BuildingResourceProductionDtoResponse `json:"productions" binding:"required"`
	Storages    []BuildingResourceStorageDtoResponse    `json:"storages" binding:"required"`
}

type BuildingCostDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid" binding:"required"`
	Cost     int       `json:"cost" binding:"required"`
	Progress float64   `json:"progress" binding:"required"`
}

type BuildingResourceProductionDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid" binding:"required"`
	Base     int       `json:"base" binding:"required"`
	Progress float64   `json:"progress" binding:"required"`
}

type BuildingResourceStorageDtoResponse struct {
	Resource uuid.UUID `json:"resource" format:"uuid" binding:"required"`
	Base     int       `json:"base" binding:"required"`
	Scale    float64   `json:"scale" binding:"required"`
	Progress float64   `json:"progress" binding:"required"`
}
