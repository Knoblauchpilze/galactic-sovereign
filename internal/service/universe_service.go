package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type UniverseService interface {
	Create(ctx context.Context, universeDto communication.UniverseDtoRequest) (communication.UniverseDtoResponse, error)
	Get(ctx context.Context, id uuid.UUID) (communication.FullUniverseDtoResponse, error)
	List(ctx context.Context) ([]communication.UniverseDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type universeServiceImpl struct {
	conn db.ConnectionPool

	universeRepo repositories.UniverseRepository
	resourceRepo repositories.ResourceRepository
	buildingRepo repositories.BuildingRepository
}

func NewUniverseService(conn db.ConnectionPool, repos repositories.Repositories) UniverseService {
	return &universeServiceImpl{
		conn:         conn,
		universeRepo: repos.Universe,
		resourceRepo: repos.Resource,
		buildingRepo: repos.Building,
	}
}

func (s *universeServiceImpl) Create(ctx context.Context, universeDto communication.UniverseDtoRequest) (communication.UniverseDtoResponse, error) {
	universe := communication.FromUniverseDtoRequest(universeDto)

	createdUniverse, err := s.universeRepo.Create(ctx, universe)
	if err != nil {
		return communication.UniverseDtoResponse{}, err
	}

	out := communication.ToUniverseDtoResponse(createdUniverse)
	return out, nil
}

func (s *universeServiceImpl) Get(ctx context.Context, id uuid.UUID) (communication.FullUniverseDtoResponse, error) {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return communication.FullUniverseDtoResponse{}, err
	}
	defer tx.Close(ctx)

	universe, err := s.universeRepo.Get(ctx, tx, id)
	if err != nil {
		return communication.FullUniverseDtoResponse{}, err
	}

	resources, err := s.resourceRepo.List(ctx, tx)
	if err != nil {
		return communication.FullUniverseDtoResponse{}, err
	}

	buildings, err := s.buildingRepo.List(ctx, tx)
	if err != nil {
		return communication.FullUniverseDtoResponse{}, err
	}

	out := communication.ToFullUniverseDtoResponse(universe, resources, buildings)

	return out, nil
}

func (s *universeServiceImpl) List(ctx context.Context) ([]communication.UniverseDtoResponse, error) {
	universes, err := s.universeRepo.List(ctx)
	if err != nil {
		return []communication.UniverseDtoResponse{}, err
	}

	var out []communication.UniverseDtoResponse
	for _, universe := range universes {
		dto := communication.ToUniverseDtoResponse(universe)
		out = append(out, dto)
	}

	return out, nil
}

func (s *universeServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	return s.universeRepo.Delete(ctx, tx, id)
}
