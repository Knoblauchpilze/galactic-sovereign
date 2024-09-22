package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/game"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type planetResourceServiceImpl struct {
	planetResourceRepo           repositories.PlanetResourceRepository
	planetResourceProductionRepo repositories.PlanetResourceProductionRepository
}

func NewPlanetResourceService(repos repositories.Repositories) game.PlanetResourceService {
	return &planetResourceServiceImpl{
		planetResourceRepo:           repos.PlanetResource,
		planetResourceProductionRepo: repos.PlanetResourceProduction,
	}
}

func (s *planetResourceServiceImpl) UpdatePlanetUntil(ctx context.Context, tx db.Transaction, planet uuid.UUID, until time.Time) error {
	resources, err := s.planetResourceRepo.ListForPlanet(ctx, tx, planet)
	if err != nil {
		return err
	}

	productions, err := s.planetResourceProductionRepo.ListForPlanet(ctx, tx, planet)
	if err != nil {
		return err
	}

	productionsMap := persistence.ToPlanetResourceProductionMap(productions)

	for _, resource := range resources {
		production, ok := productionsMap[resource.Resource]
		if !ok {
			continue
		}

		resource := game.UpdatePlanetResourceAmountToTime(resource, float64(production.Production), tx.TimeStamp())

		_, err = s.planetResourceRepo.Update(ctx, tx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}
