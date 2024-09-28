package game

import (
	"context"
	"math"
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
	PlanetResourceStorageRepo    repositories.PlanetResourceStorageRepository
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

	storages, err := data.PlanetResourceStorageRepo.ListForPlanet(ctx, tx, data.Planet)
	if err != nil {
		return err
	}

	productionsMap := toPlanetResourceProductionMap(productions)
	storagesMap := toPlanetResourceStorageMap(storages)

	for _, resource := range resources {
		production, ok := productionsMap[resource.Resource]
		if !ok {
			continue
		}

		storage := storagesMap[resource.Resource]

		resource := updatePlanetResourceAmountToTime(resource, float64(production), float64(storage), data.Until)

		_, err = data.PlanetResourceRepo.Update(ctx, tx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func toPlanetResourceProductionMap(in []persistence.PlanetResourceProduction) map[uuid.UUID]int {
	out := make(map[uuid.UUID]int)

	for _, production := range in {
		if _, ok := out[production.Resource]; !ok {
			out[production.Resource] = production.Production
		} else {
			out[production.Resource] += production.Production
		}
	}

	return out
}

func toPlanetResourceStorageMap(in []persistence.PlanetResourceStorage) map[uuid.UUID]int {
	out := make(map[uuid.UUID]int)

	for _, storage := range in {
		if _, ok := out[storage.Resource]; !ok {
			out[storage.Resource] = storage.Storage
		} else {
			out[storage.Resource] += storage.Storage
		}
	}

	return out
}

func updatePlanetResourceAmountToTime(resource persistence.PlanetResource, production float64, storage float64, moment time.Time) persistence.PlanetResource {
	elapsed := moment.Sub(resource.UpdatedAt)
	if elapsed < 0 {
		return resource
	}

	if resource.Amount < storage {
		hours := elapsed.Hours()
		resource.Amount += hours * production
		resource.Amount = math.Min(resource.Amount, storage)
	}

	resource.UpdatedAt = moment

	return resource
}
