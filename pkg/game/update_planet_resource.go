package game

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type PlanetResourceUpdateData struct {
	Planet                       uuid.UUID
	Until                        time.Time
	PlanetResourceRepo           repositories.PlanetResourceRepository
	PlanetResourceProductionRepo repositories.PlanetResourceProductionRepository
}

func UpdatePlanetResourcesToTime(ctx context.Context, tx db.Transaction, data PlanetResourceUpdateData) error {
	resources, err := data.PlanetResourceRepo.ListForPlanet(ctx, tx, data.Planet)
	if err != nil {
		return err
	}

	productions, err := data.PlanetResourceProductionRepo.ListForPlanet(ctx, tx, data.Planet)
	if err != nil {
		return err
	}

	productionsMap := persistence.ToPlanetResourceProductionMap(productions)

	for _, resource := range resources {
		production, ok := productionsMap[resource.Resource]
		if !ok {
			continue
		}

		resource := updatePlanetResourceAmountToTime(resource, float64(production.Production), data.Until)

		_, err = data.PlanetResourceRepo.Update(ctx, tx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func updatePlanetResourceAmountToTime(resource persistence.PlanetResource, production float64, moment time.Time) persistence.PlanetResource {
	elapsed := moment.Sub(resource.UpdatedAt)
	if elapsed < 0 {
		return resource
	}

	hours := elapsed.Hours()
	resource.Amount += hours * production
	resource.UpdatedAt = moment

	return resource
}
