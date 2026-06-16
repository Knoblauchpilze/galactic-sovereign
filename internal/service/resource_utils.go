package service

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
)

type updateFunc func(existing float64, change float64) float64

var addResource = func(existing float64, change float64) float64 {
	return existing + change
}
var subtractResource = func(existing float64, change float64) float64 {
	return existing - change
}

func findResourceForCost(resources []persistence.PlanetResource, cost persistence.BuildingActionCost) (persistence.PlanetResource, error) {
	for _, resource := range resources {
		if resource.Resource == cost.Resource {
			return resource, nil
		}
	}

	return persistence.PlanetResource{}, errors.FromCode(errNoSuchResource)
}

func updatePlanetResourceWithCosts(ctx context.Context, tx db.Transaction, repo repositories.PlanetResourceRepository, resources []persistence.PlanetResource, costs []persistence.BuildingActionCost, update updateFunc) error {
	for _, cost := range costs {
		resource, err := findResourceForCost(resources, cost)
		if err != nil {
			if errCode, ok := errors.AsErrorWithCode(err); ok && errCode.Code == errNoSuchResource {
				return ErrActionUsesUnknownResource
			}

			return err
		}

		planetResource := resource
		planetResource.Amount = update(planetResource.Amount, float64(cost.Amount))

		_, err = repo.Update(ctx, tx, planetResource)
		if err != nil {
			if err == repositories.ErrOptimisticLockException {
				return ErrConflictingStateForAction
			}

			return err
		}
	}

	return nil
}
