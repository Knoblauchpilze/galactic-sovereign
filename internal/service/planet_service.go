package service

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/communication"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
)

type PlanetService interface {
	Create(ctx context.Context, planetDto communication.PlanetDtoRequest) (communication.PlanetDtoResponse, error)
	Get(ctx context.Context, id uuid.UUID) (communication.FullPlanetDtoResponse, error)
	List(ctx context.Context) ([]communication.PlanetDtoResponse, error)
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]communication.PlanetDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type planetServiceImpl struct {
	conn db.Connection

	planetRepo                   repositories.PlanetRepository
	planetBuildingRepo           repositories.PlanetBuildingRepository
	planetResourceRepo           repositories.PlanetResourceRepository
	planetResourceProductionRepo repositories.PlanetResourceProductionRepository
	planetResourceStorageRepo    repositories.PlanetResourceStorageRepository
	buildingActionRepo           repositories.BuildingActionRepository
}

func NewPlanetService(conn db.Connection, repos repositories.Repositories) PlanetService {
	return &planetServiceImpl{
		conn:                         conn,
		planetRepo:                   repos.Planet,
		planetBuildingRepo:           repos.PlanetBuilding,
		planetResourceRepo:           repos.PlanetResource,
		planetResourceProductionRepo: repos.PlanetResourceProduction,
		planetResourceStorageRepo:    repos.PlanetResourceStorage,
		buildingActionRepo:           repos.BuildingAction,
	}
}

func (s *planetServiceImpl) Create(ctx context.Context, planetDto communication.PlanetDtoRequest) (communication.PlanetDtoResponse, error) {
	planet := communication.FromPlanetDtoRequest(planetDto)

	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return communication.PlanetDtoResponse{}, err
	}
	defer tx.Close(ctx)

	createdPlanet, err := s.planetRepo.Create(ctx, tx, planet)
	if err != nil {
		return communication.PlanetDtoResponse{}, err
	}

	out := communication.ToPlanetDtoResponse(createdPlanet)
	return out, nil
}

func (s *planetServiceImpl) Get(ctx context.Context, id uuid.UUID) (communication.FullPlanetDtoResponse, error) {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return communication.FullPlanetDtoResponse{}, err
	}
	defer tx.Close(ctx)

	planet, err := s.planetRepo.Get(ctx, tx, id)
	if err != nil {
		return communication.FullPlanetDtoResponse{}, err
	}

	resources, err := s.planetResourceRepo.ListForPlanet(ctx, tx, planet.Id)
	if err != nil {
		return communication.FullPlanetDtoResponse{}, err
	}

	productions, err := s.planetResourceProductionRepo.ListForPlanet(ctx, tx, planet.Id)
	if err != nil {
		return communication.FullPlanetDtoResponse{}, err
	}

	storages, err := s.planetResourceStorageRepo.ListForPlanet(ctx, tx, planet.Id)
	if err != nil {
		return communication.FullPlanetDtoResponse{}, err
	}

	buildings, err := s.planetBuildingRepo.ListForPlanet(ctx, tx, planet.Id)
	if err != nil {
		return communication.FullPlanetDtoResponse{}, err
	}

	buildingActions, err := s.buildingActionRepo.ListForPlanet(ctx, tx, planet.Id)
	if err != nil {
		return communication.FullPlanetDtoResponse{}, err
	}

	out := communication.ToFullPlanetDtoResponse(planet, resources, productions, storages, buildings, buildingActions)

	return out, nil
}

func (s *planetServiceImpl) List(ctx context.Context) ([]communication.PlanetDtoResponse, error) {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return []communication.PlanetDtoResponse{}, err
	}
	defer tx.Close(ctx)

	planets, err := s.planetRepo.List(ctx, tx)
	if err != nil {
		return []communication.PlanetDtoResponse{}, err
	}

	var out []communication.PlanetDtoResponse
	for _, planet := range planets {
		dto := communication.ToPlanetDtoResponse(planet)
		out = append(out, dto)
	}

	return out, nil
}

func (s *planetServiceImpl) ListForPlayer(ctx context.Context, player uuid.UUID) ([]communication.PlanetDtoResponse, error) {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return []communication.PlanetDtoResponse{}, err
	}
	defer tx.Close(ctx)

	planets, err := s.planetRepo.ListForPlayer(ctx, tx, player)
	if err != nil {
		return []communication.PlanetDtoResponse{}, err
	}

	var out []communication.PlanetDtoResponse
	for _, planet := range planets {
		dto := communication.ToPlanetDtoResponse(planet)
		out = append(out, dto)
	}

	return out, nil
}

func (s *planetServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	err = s.buildingActionRepo.DeleteForPlanet(ctx, tx, id)
	if err != nil {
		return err
	}

	return s.planetRepo.Delete(ctx, tx, id)
}
