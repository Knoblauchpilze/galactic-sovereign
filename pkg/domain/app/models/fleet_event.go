package models

import "time"

type FleetEvent struct {
	Timestamp time.Time
	Event     FleetEventType
}

type FleetEventType string

const (
	ArrivalEvent FleetEventType = "ARRIVAL"
	ReturnEvent  FleetEventType = "RETURN"
)
