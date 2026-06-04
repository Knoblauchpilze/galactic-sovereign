package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UniverseDtoRequest struct {
	Name string `json:"name"`
}

type UniverseDtoResponse struct {
	Id   uuid.UUID `json:"id" format:"uuid"`
	Name string    `json:"name"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
}
