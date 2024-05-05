package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ApiKeyDtoResponse struct {
	Key        uuid.UUID
	ValidUntil time.Time
}

func ToApiKeyDtoResponse(apiKey persistence.ApiKey) ApiKeyDtoResponse {
	return ApiKeyDtoResponse{
		Key:        apiKey.Key,
		ValidUntil: apiKey.ValidUntil,
	}
}
