package models

import (
	"time"

	"github.com/google/uuid"
)

type Fleet struct {
	Id     uuid.UUID
	Player uuid.UUID
	Source uuid.UUID
	// TODO: Properly handle this
	Mission     FleetMission
	Destination FleetDestination
	Ships       []FleetShip
	CreatedAt   time.Time
	ArrivalAt   time.Time
	ReturnAt    time.Time
	UpdatedAt   time.Time
	Version     int
}

type FleetDestination struct {
	Coordinate Coordinate
	Target     *uuid.UUID
}

type FleetShip struct {
	Ship  uuid.UUID
	Count int
}

func NewFleet(origin FleetOrigin, order FleetOrder, flight FleetFlight) (Fleet, error) {
	fleet := Fleet{
		Id:     uuid.New(),
		Player: origin.Player,
		Source: origin.Source,
		Destination: FleetDestination{
			Coordinate: flight.Destination,
			Target:     flight.Target,
		},
		Ships:     order.Ships,
		CreatedAt: flight.StartTime,
		ArrivalAt: flight.ArrivalTime(),
		ReturnAt:  flight.ReturnTime(),
		UpdatedAt: flight.StartTime,
		Version:   0,
	}

	return fleet, nil
}
