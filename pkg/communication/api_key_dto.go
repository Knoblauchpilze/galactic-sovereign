package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ApiKeyDtoResponse struct {
	Key        uuid.UUID `json:"key"`
	ValidUntil time.Time `json:"validUntil"`
}

func ToApiKeyDtoResponse(apiKey persistence.ApiKey) ApiKeyDtoResponse {
	return ApiKeyDtoResponse{
		Key:        apiKey.Key,
		ValidUntil: apiKey.ValidUntil,
	}
}
