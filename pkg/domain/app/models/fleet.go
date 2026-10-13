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
	Schedule    FleetFlightSchedule
	CreatedAt   time.Time
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

type FleetFlightSchedule struct {
	ArrivalAt time.Time
	ReturnAt  time.Time
	Status    FleetStatus
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
		Ships: order.Ships,
		Schedule: FleetFlightSchedule{
			ArrivalAt: flight.ArrivalTime(),
			ReturnAt:  flight.ReturnTime(),
			Status:    StatusOutbound,
		},
		CreatedAt: flight.StartTime,
		UpdatedAt: flight.StartTime,
		Version:   0,
	}

	return fleet, nil
}

func (f *Fleet) UpdateToTime(moment time.Time) error {
	if f.UpdatedAt.Compare(moment) >= 0 {
		return nil
	}

	err := f.checkUpToDate(moment)
	if err != nil {
		return err
	}

	if f.Schedule.ReturnAt.Compare(f.UpdatedAt) <= 0 {
		return domainerrors.ErrFleetAlreadyArrived
	}

	f.UpdatedAt = moment
	f.Version++

	return nil
}

func (f *Fleet) ApplyArrival(target *Planet) error {
	if (target != nil && f.Destination.Target == nil) ||
		(target == nil && f.Destination.Target != nil) ||
		(target != nil && target.Id != *f.Destination.Target) {
		return domainerrors.ErrWrongFleetTarget
	}

	if f.UpdatedAt != f.Schedule.ArrivalAt {
		return domainerrors.ErrFleetNotUpToDate
	}
	if target != nil && target.UpdatedAt != f.Schedule.ArrivalAt {
		return domainerrors.ErrPlanetNotUpToDate
	}

	f.Schedule.Status = StatusReturning
	f.Version++

	return nil
}

func (f *Fleet) ApplyReturn(target *Planet) error {
	if target == nil || target.Id != f.Source {
		return domainerrors.ErrWrongFleetTarget
	}

	if f.UpdatedAt != f.Schedule.ReturnAt {
		return domainerrors.ErrFleetNotUpToDate
	}
	if target.UpdatedAt != f.Schedule.ReturnAt {
		return domainerrors.ErrPlanetNotUpToDate
	}

	target.AddShips(f.Ships)

	return nil
}

func (f *Fleet) Completed() bool {
	return f.UpdatedAt.Compare(f.Schedule.ReturnAt) >= 0 &&
		f.Schedule.Status == StatusReturning
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

func (f *Fleet) checkUpToDate(moment time.Time) error {
	if f.Schedule.Status == StatusOutbound &&
		f.Schedule.ArrivalAt.Compare(moment) < 0 &&
		f.UpdatedAt.Compare(f.Schedule.ArrivalAt) < 0 {
		return domainerrors.ErrFleetNotUpToDate
	}

	if f.Schedule.Status == StatusReturning &&
		f.Schedule.ReturnAt.Compare(moment) < 0 &&
		f.UpdatedAt.Compare(f.Schedule.ReturnAt) < 0 {
		return domainerrors.ErrFleetNotUpToDate
	}

	return nil
}
