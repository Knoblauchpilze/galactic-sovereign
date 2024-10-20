package service

import (
	"context"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
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

	return persistence.PlanetResource{}, errors.NewCode(noSuchResource)
}

func updatePlanetResourceWithCosts(ctx context.Context, tx db.Transaction, repo repositories.PlanetResourceRepository, resources []persistence.PlanetResource, costs []persistence.BuildingActionCost, update updateFunc) error {
	for _, cost := range costs {
		resource, err := findResourceForCost(resources, cost)
		if err != nil {
			if errors.IsErrorWithCode(err, noSuchResource) {
				return errors.NewCode(FailedToCreateAction)
			}

			return err
		}

		planetResource := resource
		planetResource.Amount = update(planetResource.Amount, float64(cost.Amount))

		_, err = repo.Update(ctx, tx, planetResource)
		if err != nil {
			if errors.IsErrorWithCode(err, db.OptimisticLockException) {
				return errors.NewCode(ConflictingStateForAction)
			}

			return err
		}
	}

	return nil
}
