package request

import (
	"github.com/google/uuid"
)

type FleetCreationRequest struct {
	Planet      uuid.UUID
	Destination FleetDestinationRequest
	Ships       []FleetShipRequest
}

type FleetDestinationRequest struct {
	Galaxy      int
	SolarSystem int
	Position    int
}

type FleetShipRequest struct {
	Ship  uuid.UUID
	Count int
}
