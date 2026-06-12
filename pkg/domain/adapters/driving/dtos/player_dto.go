package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PlayerDtoRequest struct {
	ApiUser  uuid.UUID `json:"api_user" format:"uuid"`
	Universe uuid.UUID `json:"universe" format:"uuid"`
	Name     string    `json:"name"`
}

type PlayerDtoResponse struct {
	Id       uuid.UUID `json:"id" format:"uuid"`
	ApiUser  uuid.UUID `json:"api_user" format:"uuid"`
	Universe uuid.UUID `json:"universe" format:"uuid"`
	Name     string    `json:"name"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`

	Planets []uuid.UUID `json:"planets"`
}
