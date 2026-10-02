package models

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type FleetFlight struct {
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
