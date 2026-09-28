package request

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
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

func (r FleetDestinationRequest) ToCoordinates() models.Coordinate {
	return models.Coordinate{
		Galaxy:      r.Galaxy,
		SolarSystem: r.SolarSystem,
		Position:    r.Position,
	}
}
