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
	Timestamp time.Time
	Fleet     Fleet

	// Arrival: the existing destination planet, if there is one.
	// Return: the source planet.
	Planet *Planet
}

type SimulationStepResult struct {
	Planet *Planet
	Fleet  *Fleet
}
