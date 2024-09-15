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

type actionServiceImpl struct {
	conn db.ConnectionPool

	buildingActionRepo                   repositories.BuildingActionRepository
	buildingActionCostRepo               repositories.BuildingActionCostRepository
	buildingActionResourceProductionRepo repositories.BuildingActionResourceProductionRepository
	planetBuildingRepo                   repositories.PlanetBuildingRepository
	planetResourceRepo                   repositories.PlanetResourceRepository
	planetResourceProductionRepo         repositories.PlanetResourceProductionRepository
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
	}
}

func (s *actionServiceImpl) ProcessActionsUntil(ctx context.Context, until time.Time) error {
	actions, err := s.fetchActionsUntil(ctx, until)
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

func (s *actionServiceImpl) fetchActionsUntil(ctx context.Context, until time.Time) ([]persistence.BuildingAction, error) {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return []persistence.BuildingAction{}, err
	}
	defer tx.Close(ctx)

	return s.buildingActionRepo.ListBeforeCompletionTime(ctx, tx, until)
}

func (s *actionServiceImpl) processAction(ctx context.Context, action persistence.BuildingAction) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	err = s.updateResourcesForPlanet(ctx, tx, action.Planet)
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

func toPlanetResourceProductionMap(in []persistence.PlanetResourceProduction) map[uuid.UUID]persistence.PlanetResourceProduction {
	out := make(map[uuid.UUID]persistence.PlanetResourceProduction)

	for _, production := range in {
		out[production.Resource] = production
	}

	return out
}

func (s *actionServiceImpl) updateResourcesForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	resources, err := s.planetResourceRepo.ListForPlanet(ctx, tx, planet)
	if err != nil {
		return err
	}

	productions, err := s.planetResourceProductionRepo.ListForPlanet(ctx, tx, planet)
	if err != nil {
		return err
	}

	productionsMap := toPlanetResourceProductionMap(productions)

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

func toActionResourceProductionMap(in []persistence.BuildingActionResourceProduction) map[uuid.UUID]persistence.BuildingActionResourceProduction {
	out := make(map[uuid.UUID]persistence.BuildingActionResourceProduction)

	for _, production := range in {
		out[production.Resource] = production
	}

	return out
}

func (s *actionServiceImpl) updateResourcesProductionForPlanet(ctx context.Context, tx db.Transaction, action persistence.BuildingAction) error {
	production, err := s.planetResourceProductionRepo.GetForPlanetAndBuilding(ctx, tx, action.Planet, &action.Building)
	if err != nil {
		return err
	}

	newProductions, err := s.buildingActionResourceProductionRepo.ListForAction(ctx, tx, action.Id)
	if err != nil {
		return err
	}

	newProductionsMap := toActionResourceProductionMap(newProductions)

	newProduction, ok := newProductionsMap[production.Resource]
	if !ok {
		return nil
	}

	updatedProduction := persistence.ToPlanetResourceProduction(newProduction, production)
	updatedProduction.UpdatedAt = action.CompletedAt
	_, err = s.planetResourceProductionRepo.Update(ctx, tx, updatedProduction)

	return err
}
