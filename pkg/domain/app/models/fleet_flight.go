package models

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type FleetFlight struct {
	StartTime   time.Time
	Source      Coordinate
	Destination Coordinate
	Target      *uuid.UUID
	Speed       int
}

func (f FleetFlight) Duration() time.Duration {
	// https://ogame.fandom.com/wiki/Distance
	d := f.Source.DistanceTo(f.Destination)

	seconds := int(math.Floor(10.0 + 3500.0*math.Sqrt(10.0*float64(d)/float64(f.Speed))))
	return time.Duration(seconds) * time.Second
}

func (f FleetFlight) ArrivalTime() time.Time {
	return f.StartTime.Add(f.Duration())
}

func (f FleetFlight) ReturnTime() time.Time {
	return f.ArrivalTime().Add(f.Duration())
}
