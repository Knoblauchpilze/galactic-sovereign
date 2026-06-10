package request

import (
	"github.com/google/uuid"
)

type BuildingActionCreationRequest struct {
	Planet   uuid.UUID `json:"planet" format:"uuid"`
	Building uuid.UUID `json:"building" format:"uuid"`
}
