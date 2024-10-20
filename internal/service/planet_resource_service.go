package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/game"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
)

type planetResourceServiceImpl struct {
	conn db.ConnectionPool

	planetResourceRepo           repositories.PlanetResourceRepository
	planetResourceProductionRepo repositories.PlanetResourceProductionRepository
	planetResourceStorageRepo    repositories.PlanetResourceStorageRepository
}

func NewPlanetResourceService(conn db.ConnectionPool, repos repositories.Repositories) game.PlanetResourceService {
	return &planetResourceServiceImpl{
		conn: conn,

		planetResourceRepo:           repos.PlanetResource,
		planetResourceProductionRepo: repos.PlanetResourceProduction,
		planetResourceStorageRepo:    repos.PlanetResourceStorage,
	}
}

func (s *planetResourceServiceImpl) UpdatePlanetUntil(ctx context.Context, planet uuid.UUID, until time.Time) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	data := game.PlanetResourceUpdateData{
		Planet:                       planet,
		Until:                        until,
		PlanetResourceRepo:           s.planetResourceRepo,
		PlanetResourceProductionRepo: s.planetResourceProductionRepo,
		PlanetResourceStorageRepo:    s.planetResourceStorageRepo,
	}

	return game.UpdatePlanetResourcesToTime(ctx, tx, data)
}
