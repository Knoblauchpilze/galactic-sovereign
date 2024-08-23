package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/game"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type BuildingActionService interface {
	Create(ctx context.Context, actionDto communication.BuildingActionDtoRequest) (communication.BuildingActionDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type buildingActionValidator func(action persistence.BuildingAction, resources []persistence.PlanetResource, buildings []persistence.PlanetBuilding, costs []persistence.BuildingCost) error
type buildingActionConsolidator func(action persistence.BuildingAction, buildings []persistence.PlanetBuilding, resources []persistence.Resource, costs []persistence.BuildingCost) (persistence.BuildingAction, error)

type buildingActionServiceImpl struct {
	conn db.ConnectionPool

	validator    buildingActionValidator
	consolidator buildingActionConsolidator

	resourceRepo       repositories.ResourceRepository
	planetResourceRepo repositories.PlanetResourceRepository
	planetBuildingRepo repositories.PlanetBuildingRepository
	buildingCostRepo   repositories.BuildingCostRepository
	buildingActionRepo repositories.BuildingActionRepository
}

func NewBuildingActionService(conn db.ConnectionPool, repos repositories.Repositories) BuildingActionService {
	return newBuildingActionService(conn, repos, game.ValidateBuildingAction, game.ConsolidateBuildingAction)
}

func newBuildingActionService(conn db.ConnectionPool, repos repositories.Repositories, validator buildingActionValidator, consolidator buildingActionConsolidator) BuildingActionService {
	return &buildingActionServiceImpl{
		conn: conn,

		validator:    validator,
		consolidator: consolidator,

		resourceRepo:       repos.Resource,
		planetResourceRepo: repos.PlanetResource,
		planetBuildingRepo: repos.PlanetBuilding,
		buildingCostRepo:   repos.BuildingCost,
		buildingActionRepo: repos.BuildingAction,
	}
}

func (s *buildingActionServiceImpl) Create(ctx context.Context, actionDto communication.BuildingActionDtoRequest) (communication.BuildingActionDtoResponse, error) {
	action := communication.FromBuildingActionDtoRequest(actionDto)

	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}
	defer tx.Close(ctx)

	resources, err := s.resourceRepo.List(ctx, tx)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	planetResources, err := s.planetResourceRepo.ListForPlanet(ctx, tx, action.Planet)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	buildings, err := s.planetBuildingRepo.ListForPlanet(ctx, tx, action.Planet)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	costs, err := s.buildingCostRepo.ListForBuilding(ctx, tx, action.Building)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	action, err = s.consolidator(action, buildings, resources, costs)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	err = s.validator(action, planetResources, buildings, costs)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	action, err = s.buildingActionRepo.Create(ctx, tx, action)
	if err != nil {
		return communication.BuildingActionDtoResponse{}, err
	}

	out := communication.ToBuildingActionDtoResponse(action)
	return out, nil
}

func (s *buildingActionServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	return s.buildingActionRepo.Delete(ctx, tx, id)
}
