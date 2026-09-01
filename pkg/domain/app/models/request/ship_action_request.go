package request

import (
	"github.com/google/uuid"
)

type ShipActionCreationRequest struct {
	Planet uuid.UUID
	Ship   uuid.UUID
	Count  int
}
