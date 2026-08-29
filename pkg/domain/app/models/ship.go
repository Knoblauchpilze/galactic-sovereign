package models

import (
	"fmt"
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
	Resource              uuid.UUID
	Cost                  int
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
		fmt.Printf("cost: %+v, count: %v, total: %v\n", baseCost, count, cost.Amount)
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
		fmt.Printf("cost: %+v, time: %v, total: %v\n", baseCost, resourceBuildTime, buildTimeHour)
	}

	nanoSeconds := math.Floor(buildTimeHour * float64(time.Hour.Nanoseconds()))

	return time.Duration(nanoSeconds)
}
