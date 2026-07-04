package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PlayerDtoRequest struct {
	ApiUser  uuid.UUID `json:"api_user" format:"uuid" binding:"required"`
	Universe uuid.UUID `json:"universe" format:"uuid" binding:"required"`
	Name     string    `json:"name" example:"count tesla" binding:"required"`
}

type PlayerDtoResponse struct {
	Id       uuid.UUID `json:"id" format:"uuid" binding:"required"`
	ApiUser  uuid.UUID `json:"api_user" format:"uuid" binding:"required"`
	Universe uuid.UUID `json:"universe" format:"uuid" binding:"required"`
	Name     string    `json:"name" example:"emperor palpatine" binding:"required"`

	CreatedAt time.Time `json:"created_at" format:"date-time" binding:"required"`

	Homeworld uuid.UUID   `json:"homeworld" format:"uuid" binding:"required"`
	Planets   []uuid.UUID `json:"planets" format:"uuid" binding:"required"`
}
