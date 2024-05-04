package communication

import (
	"time"

	"github.com/google/uuid"
)

type ApiKeyDtoResponse struct {
	Key        uuid.UUID
	ValidUntil time.Time
}
