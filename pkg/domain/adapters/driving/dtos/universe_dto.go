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

	CreatedAt time.Time `json:"createdAt" format:"date-time" binding:"required"`

	Resources []ResourceDtoResponse `json:"resources" binding:"required"`
}

type ResourceDtoResponse struct {
	Id   uuid.UUID `json:"id" format:"uuid" binding:"required"`
	Name string    `json:"name" example:"metal" binding:"required"`

	StartAmount     int `json:"start_amount" binding:"required"`
	StartProduction int `json:"start_production" binding:"required"`
	StartStorage    int `json:"start_storage" binding:"required"`

	CreatedAt time.Time `json:"createdAt" format:"date-time" binding:"required"`
}
