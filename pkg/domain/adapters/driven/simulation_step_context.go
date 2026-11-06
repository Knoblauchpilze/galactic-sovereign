package drivenadapters

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

type simulationStepContext struct {
	planet         *models.Planet
	nextFleetEvent *models.FleetEventContext

	initialPlanetVersion int
	initialFleetVersion  int
	allowPlanetDeletion  bool

	resolver models.SimulationStepResolver
}

func newPlanetStepContext(
	planet models.Planet,
	allowPlanetDeletion bool,
	resolver models.SimulationStepResolver,
) *simulationStepContext {
	return &simulationStepContext{
		planet:               new(planet),
		initialPlanetVersion: planet.Version,
		allowPlanetDeletion:  allowPlanetDeletion,
		resolver:             resolver,
	}
}

func newFleetStepContext(
	nextFleetEvent models.FleetEventContext,
	resolver models.SimulationStepResolver,
) *simulationStepContext {
	ctx := simulationStepContext{
		nextFleetEvent:      new(nextFleetEvent),
		initialFleetVersion: nextFleetEvent.Fleet.Version,
		// Do not allow planet deletion for fleet event
		allowPlanetDeletion: false,
		resolver:            resolver,
	}

	if nextFleetEvent.Planet != nil {
		ctx.initialPlanetVersion = nextFleetEvent.Planet.Version
	}

	return &ctx
}

func (c *simulationStepContext) resolve(
	ctx context.Context,
	tx database.Transaction,
	timestamp time.Time,
) (models.SimulationStepResult, error) {
	step := models.SimulationStep{
		Timestamp:      timestamp,
		Planet:         c.planet,
		NextFleetEvent: c.nextFleetEvent,
	}

	result, err := c.resolver(step)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	// Reject deletion when the context does not allow it
	if !c.allowPlanetDeletion && result.Resolution == models.PlanetDeleted {
		return models.SimulationStepResult{}, domainerrors.ErrInvalidSimulationStepResult
	}
	// Reject planet update when fleet event is available
	if c.nextFleetEvent != nil && result.Resolution == models.PlanetUpdated {
		return models.SimulationStepResult{}, domainerrors.ErrInvalidSimulationStepResult
	}
	// Reject fleet creation and planet update when planet deletion is allowed
	if c.nextFleetEvent == nil && c.allowPlanetDeletion && result.Resolution != models.PlanetDeleted {
		return models.SimulationStepResult{}, domainerrors.ErrInvalidSimulationStepResult
	}
	// Reject fleet deletion when fleet event is not available
	if c.nextFleetEvent == nil && result.Resolution == models.FleetDeleted {
		return models.SimulationStepResult{}, domainerrors.ErrInvalidSimulationStepResult
	}

	planet, err := c.persistPlanet(ctx, tx, result)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	fleet, err := c.persistFleet(ctx, tx, result)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	return models.SimulationStepResult{
		Resolution: result.Resolution,
		Planet:     planet,
		Fleet:      fleet,
	}, nil
}

func (c *simulationStepContext) persistPlanet(
	ctx context.Context,
	tx database.Transaction,
	result models.SimulationStepResult,
) (*models.Planet, error) {
	if result.Resolution == models.PlanetDeleted {
		err := deletePlanetAndDetails(ctx, tx, c.planet.Id)
		if err != nil {
			return nil, err
		}
	}

	if result.Planet == nil {
		return nil, nil
	}

	out, err := saveAndReloadPlanet(ctx, tx, *result.Planet, c.initialPlanetVersion)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *simulationStepContext) persistFleet(
	ctx context.Context,
	tx database.Transaction,
	result models.SimulationStepResult,
) (*models.Fleet, error) {
	if result.Resolution == models.FleetDeleted {
		err := deleteFleetAndDetails(ctx, tx, c.nextFleetEvent.Fleet.Id)
		if err != nil {
			return nil, err
		}
	}

	if result.Fleet == nil {
		return nil, nil
	}

	if result.Resolution == models.FleetCreated {
		err := createFleetWithDetails(ctx, tx, *result.Fleet)
		if err != nil {
			if isFleetMissingTargetPlanetErr(err) {
				return nil, domainerrors.ErrFleetDestinationInvalid
			}
			if isFleetSourcePlanetErr(err) {
				return nil, domainerrors.ErrPlanetNotFound
			}
			return nil, err
		}

		return result.Fleet, nil
	}

	out, err := saveAndReloadFleet(ctx, tx, *result.Fleet, c.initialFleetVersion)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
