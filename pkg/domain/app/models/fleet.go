package models

import (
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
)

type Fleet struct {
	Id          uuid.UUID
	Player      uuid.UUID
	Source      uuid.UUID
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
	if len(order.Ships) == 0 {
		return Fleet{}, domainerrors.ErrNoShipInFleet
	}

	if err := validateMission(order, flight); err != nil {
		return Fleet{}, err
	}

	fleet := Fleet{
		Id:      uuid.New(),
		Player:  origin.Player,
		Source:  origin.Source,
		Mission: order.Mission,
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

func validateMission(order FleetOrder, flight FleetFlight) error {
	switch order.Mission {
	case MissionColonize:
		return validateExplorationMission(flight)
	case MissionTransport:
		return validateDirectedMission(flight)
	default:
		return domainerrors.ErrInvalidFleetConfiguration
	}
}

func validateExplorationMission(flight FleetFlight) error {
	if !flight.HasDistinctSourceAndDestination() || flight.Target != nil {
		return domainerrors.ErrInvalidFleetConfiguration
	}

	return nil
}

func validateDirectedMission(flight FleetFlight) error {
	if !flight.HasDistinctSourceAndDestination() || flight.Target == nil {
		return domainerrors.ErrInvalidFleetConfiguration
	}

	return nil
}
