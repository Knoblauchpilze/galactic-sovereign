package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/game"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type planetResourceServiceImpl struct {
	conn db.ConnectionPool

	planetResourceRepo           repositories.PlanetResourceRepository
	planetResourceProductionRepo repositories.PlanetResourceProductionRepository
}

func NewPlanetResourceService(conn db.ConnectionPool, repos repositories.Repositories) game.PlanetResourceService {
	return &planetResourceServiceImpl{
		conn: conn,

		planetResourceRepo:           repos.PlanetResource,
		planetResourceProductionRepo: repos.PlanetResourceProduction,
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
	}

	return game.UpdatePlanetResourcesToTime(ctx, tx, data)
}
