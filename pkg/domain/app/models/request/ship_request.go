package request

import (
	"github.com/google/uuid"
)

type ShipCreationRequest struct {
	Planet uuid.UUID
	Ship   uuid.UUID
	Count  int
}
