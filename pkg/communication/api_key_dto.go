package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ApiKeyDtoResponse struct {
	User       uuid.UUID `json:"user"`
	Key        uuid.UUID `json:"key"`
	ValidUntil time.Time `json:"validUntil"`
}

func ToApiKeyDtoResponse(apiKey persistence.ApiKey) ApiKeyDtoResponse {
	return ApiKeyDtoResponse{
		User:       apiKey.ApiUser,
		Key:        apiKey.Key,
		ValidUntil: apiKey.ValidUntil,
	}
}
