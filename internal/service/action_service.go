package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/game"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
)

type actionServiceImpl struct {
	conn db.ConnectionPool

	buildingActionRepo                   repositories.BuildingActionRepository
	buildingActionCostRepo               repositories.BuildingActionCostRepository
	buildingActionResourceProductionRepo repositories.BuildingActionResourceProductionRepository
	planetBuildingRepo                   repositories.PlanetBuildingRepository
	planetResourceRepo                   repositories.PlanetResourceRepository
	planetResourceProductionRepo         repositories.PlanetResourceProductionRepository
	planetResourceStorageRepo            repositories.PlanetResourceStorageRepository
}

func NewActionService(conn db.ConnectionPool, repos repositories.Repositories) game.ActionService {
	return &actionServiceImpl{
		conn: conn,

		buildingActionRepo:                   repos.BuildingAction,
		buildingActionCostRepo:               repos.BuildingActionCost,
		buildingActionResourceProductionRepo: repos.BuildingActionResourceProduction,
		planetBuildingRepo:                   repos.PlanetBuilding,
		planetResourceRepo:                   repos.PlanetResource,
		planetResourceProductionRepo:         repos.PlanetResourceProduction,
		planetResourceStorageRepo:            repos.PlanetResourceStorage,
	}
}

func (s *actionServiceImpl) ProcessActionsUntil(ctx context.Context, planet uuid.UUID, until time.Time) error {
	actions, err := s.fetchActionsUntil(ctx, planet, until)
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := s.processAction(ctx, action)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *actionServiceImpl) fetchActionsUntil(ctx context.Context, planet uuid.UUID, until time.Time) ([]persistence.BuildingAction, error) {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return []persistence.BuildingAction{}, err
	}
	defer tx.Close(ctx)

	return s.buildingActionRepo.ListBeforeCompletionTime(ctx, tx, planet, until)
}

func (s *actionServiceImpl) processAction(ctx context.Context, action persistence.BuildingAction) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	data := game.PlanetResourceUpdateData{
		Planet:                       action.Planet,
		Until:                        action.CompletedAt,
		PlanetResourceRepo:           s.planetResourceRepo,
		PlanetResourceProductionRepo: s.planetResourceProductionRepo,
		PlanetResourceStorageRepo:    s.planetResourceStorageRepo,
	}
	err = game.UpdatePlanetResourcesToTime(ctx, tx, data)
	if err != nil {
		return err
	}

	building, err := s.planetBuildingRepo.GetForPlanetAndBuilding(ctx, tx, action.Planet, action.Building)
	if err != nil {
		return err
	}

	updatedBuilding := persistence.ToPlanetBuilding(action, building)

	_, err = s.planetBuildingRepo.Update(ctx, tx, updatedBuilding)
	if err != nil {
		return err
	}

	err = s.updateResourcesProductionForPlanet(ctx, tx, action)
	if err != nil {
		return err
	}

	err = s.buildingActionResourceProductionRepo.DeleteForAction(ctx, tx, action.Id)
	if err != nil {
		return err
	}

	err = s.buildingActionCostRepo.DeleteForAction(ctx, tx, action.Id)
	if err != nil {
		return err
	}

	return s.buildingActionRepo.Delete(ctx, tx, action.Id)
}

func (s *actionServiceImpl) updateResourcesProductionForPlanet(ctx context.Context, tx db.Transaction, action persistence.BuildingAction) error {
	newProductions, err := s.buildingActionResourceProductionRepo.ListForAction(ctx, tx, action.Id)
	if err != nil {
		return err
	}

	for _, newProduction := range newProductions {
		err = s.updatePlanetProductionForResource(
			ctx, tx, action, newProduction)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *actionServiceImpl) updatePlanetProductionForResource(
	ctx context.Context,
	tx db.Transaction,
	action persistence.BuildingAction,
	newProduction persistence.BuildingActionResourceProduction) error {

	production, err := s.planetResourceProductionRepo.GetForPlanetAndBuilding(ctx, tx, action.Planet, &action.Building)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return s.createPlanetProductionForResourceAndBuilding(ctx, tx, action, newProduction)
		}

		return err
	}

	updatedProduction := persistence.MergeWithPlanetResourceProduction(newProduction, production)
	updatedProduction.UpdatedAt = action.CompletedAt
	_, err = s.planetResourceProductionRepo.Update(ctx, tx, updatedProduction)

	return err
}

func (s *actionServiceImpl) createPlanetProductionForResourceAndBuilding(ctx context.Context,
	tx db.Transaction,
	action persistence.BuildingAction,
	newProduction persistence.BuildingActionResourceProduction) error {
	planetProduction := persistence.ToPlanetResourceProduction(newProduction, action)

	_, err := s.planetResourceProductionRepo.Create(ctx, tx, planetProduction)

	return err
}
