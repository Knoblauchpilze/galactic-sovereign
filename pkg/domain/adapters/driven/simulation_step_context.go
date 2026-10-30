package drivenadapters

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

type simulationStepContext struct {
	initialPlanetVersion int
	initialFleetVersion  int

	step models.SimulationStep

	resolver models.SimulationStepResolver

	result models.SimulationStepResult
}

func newSimulationStepContext(
	timestamp time.Time,
	planet *models.Planet,
	nextFleetEvent *models.FleetEventContext,
	resolver models.SimulationStepResolver,
) *simulationStepContext {
	ctx := simulationStepContext{
		step: models.SimulationStep{
			Timestamp:      timestamp,
			Planet:         planet,
			NextFleetEvent: nextFleetEvent,
		},
		resolver: resolver,
	}

	if planet != nil {
		ctx.initialPlanetVersion = planet.Version
	}

	if nextFleetEvent != nil {
		if nextFleetEvent.Planet != nil {
			ctx.initialPlanetVersion = nextFleetEvent.Planet.Version
		}
		ctx.initialFleetVersion = nextFleetEvent.Fleet.Version
	}

	return &ctx
}

func (c *simulationStepContext) resolve() error {
	result, err := c.resolver(c.step)
	if err != nil {
		return err
	}

	c.result = result

	return nil
}

func (c *simulationStepContext) persist(
	ctx context.Context,
	tx database.Transaction,
) (models.SimulationStepResult, error) {
	planet, err := c.persistPlanet(ctx, tx)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	fleet, err := c.persistFleet(ctx, tx)
	if err != nil {
		return models.SimulationStepResult{}, err
	}

	return models.SimulationStepResult{
		Resolution: c.result.Resolution,
		Planet:     planet,
		Fleet:      fleet,
	}, nil
}

func (c *simulationStepContext) persistPlanet(ctx context.Context, tx database.Transaction) (*models.Planet, error) {
	if c.result.Planet == nil {
		return nil, nil
	}

	out, err := saveAndReloadPlanet(ctx, tx, *c.result.Planet, c.initialPlanetVersion)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *simulationStepContext) persistFleet(ctx context.Context, tx database.Transaction) (*models.Fleet, error) {
	if c.result.Resolution == models.FleetDeleted {
		err := deleteFleetAndDetails(ctx, tx, c.step.NextFleetEvent.Fleet.Id)
		if err != nil {
			return nil, err
		}
	}

	if c.result.Fleet == nil {
		return nil, nil
	}

	if c.result.Resolution == models.FleetCreated {
		err := createFleetWithDetails(ctx, tx, *c.result.Fleet)
		if err != nil {
			if isFleetMissingTargetPlanetErr(err) {
				return nil, domainerrors.ErrFleetDestinationInvalid
			}
			if isFleetSourcePlanetErr(err) {
				return nil, domainerrors.ErrPlanetNotFound
			}
			return nil, err
		}

		return c.result.Fleet, nil
	}

	out, err := saveAndReloadFleet(ctx, tx, *c.result.Fleet, c.initialFleetVersion)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
