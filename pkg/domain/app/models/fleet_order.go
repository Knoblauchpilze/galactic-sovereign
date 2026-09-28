package models

type FleetOrder struct {
	Mission     FleetMission
	Destination Coordinate
	Ships       []FleetShip
}
