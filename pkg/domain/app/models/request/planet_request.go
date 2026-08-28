package request

import (
	"github.com/google/uuid"
)

type PlanetCreationRequest struct {
	Player uuid.UUID
}
