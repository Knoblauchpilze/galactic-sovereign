package request

import (
	"github.com/google/uuid"
)

type FleetCreationRequest struct {
	Planet uuid.UUID
	Ships  []FleetShipRequest
}

type FleetShipRequest struct {
	Ship  uuid.UUID
	Count int
}
