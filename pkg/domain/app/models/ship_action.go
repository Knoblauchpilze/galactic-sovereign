package models

import (
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
)

type ShipAction struct {
	Id   uuid.UUID
	Ship uuid.UUID

	Count int

	// CreatedAt represents the time at which the action will start being worked
	// on by the shipyard of the planet. This can either correspond to the time
	// of creation in case there are no action currently running on the planet
	// or to the earliest time where all actions are finished.
	CreatedAt time.Time

	// NextCompletionAt represents the next timestamp at which a unit of the ship
	// produced by the action will be available in the planet's fleet. Calculating
	// the difference between CreatedAt and this field gives the time it takes to
	// produce a single ship.
	NextCompletionAt time.Time

	// CompletedAt represents the time at which the last ship of this action will
	// be finished and delivered as an element of the planet's fleet. It's also
	// the earliest availability for another ship action to start being processed
	// on the planet the action belongs to.
	CompletedAt time.Time

	// Costs are currently only available when the action is first created. When
	// it is persisted by the adapter, the costs are lost because they do not
	// play a relevant role once the action has been created as it cannot be
	// cancelled.
	Costs []ShipActionCost
}

type ShipActionCost struct {
	Resource uuid.UUID
	Amount   int
}

// TODO: Add tests for this
func (a *ShipAction) CompleteOne() error {
	if a.Count == 0 {
		return domainerrors.ErrShipActionAlreadyCompleted
	}

	// TODO: Should also bump the NextCompletionAt
	// TODO: The problem is that there's no way to no the individual completion time
	// Probably CompletedAt could be removed and repurposed to hold the individual
	// completion time
	a.Count--
	return nil
}
