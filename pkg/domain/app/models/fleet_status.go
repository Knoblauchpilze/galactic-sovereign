package models

type FleetStatus string

const (
	StatusOutbound  FleetStatus = "OUTBOUND"
	StatusReturning FleetStatus = "RETURNING"
)
