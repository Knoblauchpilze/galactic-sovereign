package models

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type Ship struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time

	Costs []ShipCost
}

type ShipCost struct {
	Resource uuid.UUID
	Cost     int

	// Defines how long each unit of a resource takes to be transformed into
	// a ship. A value of 1 means that if a ship costs 3 unit of the resource
	// it will take 3 hours to be built.
	BuildTimeHoursPerUnit float64
}

func (s Ship) CreateShipAction(
	count int,
	createdAt time.Time,
) ShipAction {
	costs := s.determineActionCost(count)

	buildTime := s.determineBuildTime()
	completionTime := createdAt.Add(time.Duration(count) * buildTime)

	action := ShipAction{
		Id:    uuid.New(),
		Ship:  s.Id,
		Count: count,

		CreatedAt:        createdAt,
		NextCompletionAt: createdAt.Add(buildTime),
		CompletedAt:      completionTime,

		Costs: costs,
	}
	return action
}

func (s Ship) determineActionCost(
	count int,
) []ShipActionCost {
	costs := []ShipActionCost{}

	for _, baseCost := range s.Costs {
		cost := ShipActionCost{
			Resource: baseCost.Resource,
			Amount:   baseCost.Cost * count,
		}
		costs = append(costs, cost)
	}

	return costs
}

func (s Ship) determineBuildTime() time.Duration {
	temp := make(map[uuid.UUID]ShipCost)
	for _, cost := range s.Costs {
		temp[cost.Resource] = cost
	}

	buildTimeHour := 0.0
	for _, baseCost := range s.Costs {
		resourceBuildTime := float64(baseCost.Cost) * baseCost.BuildTimeHoursPerUnit
		buildTimeHour += resourceBuildTime
	}

	nanoSeconds := math.Floor(buildTimeHour * float64(time.Hour.Nanoseconds()))

	return time.Duration(nanoSeconds)
}
