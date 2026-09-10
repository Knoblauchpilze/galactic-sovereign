package models

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type Shipyard struct {
	// throughput defines how fast a ship work amount is processed by the
	// shipyard. A ship lasting 1h to build can be built faster when the
	// throughput is higher than 1. The througput is a function of the
	// level of buildings that affect ship production on the planet (such
	// as the shipyard).
	throughput float64
}

func (s Shipyard) CreateShipAction(ship Ship, count int, startedAt time.Time) ShipAction {
	return ShipAction{
		Id:                 uuid.New(),
		Ship:               ship.Id,
		Count:              count,
		CreatedAt:          startedAt,
		NextCompletionAt:   startedAt.Add(s.durationFor(ship)),
		UnitCompletionTime: s.durationFor(ship),
		Costs:              ship.CostsFor(count),
	}
}

func (s Shipyard) durationFor(ship Ship) time.Duration {
	nanoseconds := float64(ship.BuildTime().Nanoseconds()) / s.throughput
	return time.Duration(math.Floor(nanoseconds))
}
