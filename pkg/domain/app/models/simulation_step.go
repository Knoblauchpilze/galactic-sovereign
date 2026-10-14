package models

import (
	"time"
)

type SimulationStep struct {
	Timestamp time.Time

	Planet Planet

	// Earliest unresolved fleet event in the requested scope.
	// Includes its fleet and affected domain state.
	NextFleetEvent *FleetEventContext
}

type FleetEventContext struct {
	// Timestamp: the timestamp of the event. It corresponds to the moment
	// the event takes place. It can differ from the parent simulation step's
	// timestamp: it is always at most equal to the parent timestamp.
	Timestamp time.Time

	Fleet Fleet

	// Arrival: the existing destination planet, if there is one.
	// Return: the source planet.
	Planet *Planet
}

type ResolutionResult int

const (
	FleetEventResolved ResolutionResult = iota
	FleetDeleted
	PlanetUpdated
)

type SimulationStepResult struct {
	Resolution ResolutionResult
	Planet     *Planet
	Fleet      *Fleet
}
