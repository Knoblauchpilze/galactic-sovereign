package service

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/game"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
)

type BuildingActionService interface {
	Create(ctx context.Context, actionDto communication.BuildingActionDtoRequest) (communication.BuildingActionDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type buildingActionCompletionTimeConsolidator func(action persistence.BuildingAction, resources []persistence.Resource, costs []persistence.BuildingActionCost) (persistence.BuildingAction, error)
type buildingActionValidator func(action persistence.BuildingAction, resources []persistence.PlanetResource, buildings []persistence.PlanetBuilding, costs []persistence.BuildingActionCost) error

type buildingActionServiceImpl struct {
	conn db.Connection

	completionTimeConsolidator buildingActionCompletionTimeConsolidator
	validator                  buildingActionValidator

	resourceRepo                         repositories.ResourceRepository
	planetResourceRepo                   repositories.PlanetResourceRepository
	planetBuildingRepo                   repositories.PlanetBuildingRepository
	buildingCostRepo                     repositories.BuildingCostRepository
	buildingResourceProduction           repositories.BuildingResourceProductionRepository
	buildingResourceStorage              repositories.BuildingResourceStorageRepository
	buildingActionRepo                   repositories.BuildingActionRepository
	buildingActionCostRepo               repositories.BuildingActionCostRepository
	buildingActionResourceProductionRepo repositories.BuildingActionResourceProductionRepository
	buildingActionResourceStorageRepo    repositories.BuildingActionResourceStorageRepository
}

func NewBuildingActionService(conn db.Connection, repos repositories.Repositories) BuildingActionService {
	return newBuildingActionServiceWithCompletionTime(conn, repos, game.ConsolidateBuildingActionCompletionTime)
}

func newBuildingActionServiceWithCompletionTime(conn db.Connection, repos repositories.Repositories, completionTimeConsolidator buildingActionCompletionTimeConsolidator) BuildingActionService {
	return newBuildingActionService(conn, repos, completionTimeConsolidator, game.ValidateBuildingAction)
}

func newBuildingActionService(
	conn db.Connection,
	repos repositories.Repositories,
	completionTimeConsolidator buildingActionCompletionTimeConsolidator,
	validator buildingActionValidator) BuildingActionService {
	return &buildingActionServiceImpl{
		conn: conn,

		completionTimeConsolidator: completionTimeConsolidator,
		validator:                  validator,

		resourceRepo:                         repos.Resource,
		planetResourceRepo:                   repos.PlanetResource,
		planetBuildingRepo:                   repos.PlanetBuilding,
		buildingCostRepo:                     repos.BuildingCost,
		buildingResourceProduction:           repos.BuildingResourceProduction,
		buildingResourceStorage:              repos.BuildingResourceStorage,
		buildingActionRepo:                   repos.BuildingAction,
		buildingActionCostRepo:               repos.BuildingActionCost,
		buildingActionResourceProductionRepo: repos.BuildingActionResourceProduction,
		buildingActionResourceStorageRepo:    repos.BuildingActionResourceStorage,
	}
}

func (s *buildingActionServiceImpl) Create(
	ctx context.Context,
	actionDto communication.BuildingActionDtoRequest,
) (communication.BuildingActionDtoResponse, error) {
	action := communication.FromBuildingActionDtoRequest(actionDto)

	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}
	defer tx.Close(ctx)

	planetResources, err := s.planetResourceRepo.ListForPlanet(ctx, tx, action.Planet)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	action, costs, err := s.consolidateAction(ctx, tx, action, planetResources)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	productions, err := s.determineBuildingActionResourceProductions(ctx, tx, action)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	storages, err := s.determineBuildingActionResourceStorages(ctx, tx, action)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	action, err = s.createAction(ctx, tx, action, costs, productions, storages, planetResources)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	out := communication.ToBuildingActionDtoResponse(action)
	return out, nil
}

func (s *buildingActionServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	action, err := s.buildingActionRepo.Get(ctx, tx, id)
	if err != nil {
		return err
	}

	if action.CompletedAt.Before(tx.TimeStamp()) {
		return errors.NewCode(ActionAlreadyCompleted)
	}

	costs, err := s.buildingActionCostRepo.ListForAction(ctx, tx, action.Id)
	if err != nil {
		return err
	}

	planetResources, err := s.planetResourceRepo.ListForPlanet(ctx, tx, action.Planet)
	if err != nil {
		return err
	}

	err = updatePlanetResourceWithCosts(ctx, tx, s.planetResourceRepo, planetResources, costs, addResource)
	if err != nil {
		return err
	}

	return s.buildingActionRepo.Delete(ctx, tx, id)
}

func (s *buildingActionServiceImpl) consolidateAction(
	ctx context.Context,
	tx db.Transaction,
	action persistence.BuildingAction,
	planetResources []persistence.PlanetResource,
) (persistence.BuildingAction, []persistence.BuildingActionCost, error) {
	var costs []persistence.BuildingActionCost

	resources, err := s.resourceRepo.List(ctx, tx)
	if err != nil {
		return action, costs, err
	}

	buildings, err := s.planetBuildingRepo.ListForPlanet(ctx, tx, action.Planet)
	if err != nil {
		return action, costs, err
	}

	action = game.ConsolidateBuildingActionLevel(action, buildings)

	costs, err = s.determineBuildingActionCosts(ctx, tx, action)
	if err != nil {
		return action, costs, err
	}

	action, err = s.completionTimeConsolidator(action, resources, costs)
	if err != nil {
		return action, costs, err
	}

	err = s.validator(action, planetResources, buildings, costs)

	return action, costs, err
}

func (s *buildingActionServiceImpl) determineBuildingActionCosts(
	ctx context.Context,
	tx db.Transaction,
	action persistence.BuildingAction,
) ([]persistence.BuildingActionCost, error) {
	var costs []persistence.BuildingActionCost

	baseCosts, err := s.buildingCostRepo.ListForBuilding(ctx, tx, action.Building)
	if err != nil {
		return costs, err
	}

	costs = game.DetermineBuildingActionCost(action, baseCosts)

	return costs, nil
}

func (s *buildingActionServiceImpl) determineBuildingActionResourceProductions(
	ctx context.Context,
	tx db.Transaction,
	action persistence.BuildingAction,
) ([]persistence.BuildingActionResourceProduction, error) {
	var productions []persistence.BuildingActionResourceProduction

	baseProductions, err := s.buildingResourceProduction.ListForBuilding(ctx, tx, action.Building)
	if err != nil {
		return productions, err
	}

	productions = game.DetermineBuildingActionResourceProduction(action, baseProductions)

	return productions, nil
}

func (s *buildingActionServiceImpl) determineBuildingActionResourceStorages(
	ctx context.Context,
	tx db.Transaction,
	action persistence.BuildingAction,
) ([]persistence.BuildingActionResourceStorage, error) {
	var storages []persistence.BuildingActionResourceStorage

	baseStorages, err := s.buildingResourceStorage.ListForBuilding(ctx, tx, action.Building)
	if err != nil {
		return storages, err
	}

	storages = game.DetermineBuildingActionResourceStorage(action, baseStorages)

	return storages, nil
}

func (s *buildingActionServiceImpl) createAction(ctx context.Context,
	tx db.Transaction,
	action persistence.BuildingAction,
	costs []persistence.BuildingActionCost,
	productions []persistence.BuildingActionResourceProduction,
	storages []persistence.BuildingActionResourceStorage,
	planetResources []persistence.PlanetResource,
) (persistence.BuildingAction, error) {
	err := updatePlanetResourceWithCosts(ctx, tx, s.planetResourceRepo, planetResources, costs, subtractResource)
	if err != nil {
		return action, err
	}

	action, err = s.buildingActionRepo.Create(ctx, tx, action)
	if err != nil {
		return action, err
	}

	for _, cost := range costs {
		_, err = s.buildingActionCostRepo.Create(ctx, tx, cost)
		if err != nil {
			return action, err
		}
	}

	for _, production := range productions {
		_, err = s.buildingActionResourceProductionRepo.Create(ctx, tx, production)
		if err != nil {
			return action, err
		}
	}

	for _, storage := range storages {
		_, err = s.buildingActionResourceStorageRepo.Create(ctx, tx, storage)
		if err != nil {
			return action, err
		}
	}

	return action, nil
}
