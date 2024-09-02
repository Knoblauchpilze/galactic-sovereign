package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/game"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
)

type actionServiceImpl struct {
	conn db.ConnectionPool

	buildingActionRepo     repositories.BuildingActionRepository
	buildingActionCostRepo repositories.BuildingActionCostRepository
	planetBuildingRepo     repositories.PlanetBuildingRepository
}

func NewActionService(conn db.ConnectionPool, repos repositories.Repositories) game.ActionService {
	return &actionServiceImpl{
		conn: conn,

		buildingActionRepo:     repos.BuildingAction,
		buildingActionCostRepo: repos.BuildingActionCost,
		planetBuildingRepo:     repos.PlanetBuilding,
	}
}

func (s *actionServiceImpl) ProcessActionsUntil(ctx context.Context, until time.Time) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	actions, err := s.buildingActionRepo.ListBeforeCompletionTime(ctx, tx, until)
	if err != nil {
		return err
	}

	for _, action := range actions {
		building, err := s.planetBuildingRepo.GetForPlanetAndBuilding(ctx, tx, action.Planet, action.Building)
		if err != nil {
			return err
		}

		updatedBuilding := persistence.ToPlanetBuilding(action, building)

		_, err = s.planetBuildingRepo.Update(ctx, tx, updatedBuilding)
		if err != nil {
			return err
		}

		err = s.buildingActionCostRepo.DeleteForAction(ctx, tx, action.Id)
		if err != nil {
			return err
		}

		err = s.buildingActionRepo.Delete(ctx, tx, action.Id)
		if err != nil {
			return err
		}
	}

	return nil
}
