package request

import (
	"github.com/google/uuid"
)

type BuildingActionCreationRequest struct {
	Planet   uuid.UUID
	Building uuid.UUID
}
