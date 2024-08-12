package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/game"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
)

type BuildingActionService interface {
	Create(ctx context.Context, actionDto communication.BuildingActionDtoRequest) (communication.BuildingActionDtoResponse, error)
}

type buildingActionValidator func(action persistence.BuildingAction, resources []persistence.PlanetResource, costs []persistence.BuildingCost, buildings []persistence.PlanetBuilding) error

type buildingActionServiceImpl struct {
	conn db.ConnectionPool

	validator buildingActionValidator

	planetResourceRepo repositories.PlanetResourceRepository
	planetBuildingRepo repositories.PlanetBuildingRepository
	buildingCostRepo   repositories.BuildingCostRepository
	buildingActionRepo repositories.BuildingActionRepository
}

func NewBuildingActionService(conn db.ConnectionPool, repos repositories.Repositories) BuildingActionService {
	return newBuildingActionService(conn, repos, game.ValidateBuildingAction)
}

func newBuildingActionService(conn db.ConnectionPool, repos repositories.Repositories, validator buildingActionValidator) BuildingActionService {
	return &buildingActionServiceImpl{
		conn:               conn,
		validator:          validator,
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

	resources, err := s.planetResourceRepo.ListForPlanet(ctx, tx, action.Planet)
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

	err = s.validator(action, resources, costs, buildings)
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
